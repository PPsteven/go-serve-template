// Copy From gin-gonic/gin/logger.go
// Aim to support add more custom params in log
// Alternative: https://github.com/gin-contrib/zap

package middleware

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-isatty"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/willf/pad"
	//"github.com/go-eagle/eagle/pkg/app"
	//"github.com/go-eagle/eagle/pkg/errcode"
	//"github.com/go-eagle/eagle/pkg/log"
)

type consoleColorModeValue int

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var consoleColorMode = autoColor

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter interface {
	BeforeResponse(c *gin.Context)

	AfterResponse(c *gin.Context)

	Format(c *gin.Context) string

	i()
}

var _ LogFormatter = (*defaultLogFormatter)(nil)

// defaultLogFormatter is the default log format function Logger middleware uses.
type defaultLogFormatter struct{}

func (f *defaultLogFormatter) BeforeResponse(c *gin.Context) {}

func (f *defaultLogFormatter) AfterResponse(c *gin.Context) {}

func (f *defaultLogFormatter) Format(c *gin.Context) string {
	var statusColor, methodColor, resetColor string
	var param = basicParams

	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s | %s %s |%s %-7s %s %#v\n%s",
		param.StartTime.Format(time.DateTime),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		GetRequestIDFromContext(c),
		"",
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func (f *defaultLogFormatter) i() {}

var basicParams *LogBasicParams

// LogBasicParams is the structure any formatter will be handed when time to log comes
type LogBasicParams struct {
	Request *http.Request

	// StartTime shows the time of the request.
	StartTime time.Time
	// EndTime shows the time request is finished.
	EndTime time.Time
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// StatusCode is HTTP response code.
	StatusCode int
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]any
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogBasicParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogBasicParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogBasicParams) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (p *LogBasicParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.isTerm)
}

// defaultLogFormatter is the default log format function Logger middleware uses.
//var defaultLogFormatter = func(param LogBasicParams) string {
//}

// DisableConsoleColor disables color output in the console.
func DisableConsoleColor() {
	consoleColorMode = disableColor
}

// ForceConsoleColor force color output in the console.
func ForceConsoleColor() {
	consoleColorMode = forceColor
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger is a middleware function that logs the each request.
func Logger() gin.HandlerFunc { return LoggerWithConfig(LoggerConfig{}) }

// LoggerWithConfig is same as Logger() but with custom config.
func LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = &defaultLogFormatter{}
	}

	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	notlogged := conf.SkipPaths

	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		if raw != "" {
			path = path + "?" + raw
		}

		// Read the Body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		//log.Debugf("New request come in, path: %s, Method: %s, body `%s`", path, method, string(bodyBytes))
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		formatter.BeforeResponse(c)

		//defer func() {
		//
		//	var code int
		//	var message string
		//
		//	// get code and message
		//	var response app.Response
		//	if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
		//		log.Errorf("response body can not unmarshal to model.Response struct, body: `%s`, err: %+v",
		//			blw.body.Bytes(), err)
		//		code = errcode.ErrInternalServer.Code()
		//		message = err.Error()
		//	} else {
		//		code = response.Code
		//		message = response.Message
		//	}
		//}()

		// Continue.
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Calculates the latency.
			end := time.Now().UTC()

			basicParams = &LogBasicParams{
				Request: c.Request,
				isTerm:  isTerm,
				Keys:    c.Keys,
				// Response Time
				StartTime: start,
				EndTime:   end,
				Latency:   end.Sub(start),
				// Detail
				ClientIP:     c.ClientIP(),
				Method:       c.Request.Method,
				StatusCode:   c.Writer.Status(),
				ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
				BodySize:     c.Writer.Size(),
				Path:         path,
			}

			formatter.AfterResponse(c)

			fmt.Fprint(out, formatter.Format(c))
		}
	}
}
