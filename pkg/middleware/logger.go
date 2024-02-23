// Copy From gin-gonic/gin/logger.go
// Aim to support add more custom params in log
// Alternative: https://github.com/gin-contrib/zap

package middleware

import (
	"bytes"
	"fmt"
	"go-server-template/pkg/logger"
	"go-server-template/pkg/trace"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string

	Filter func(ctx *gin.Context) bool

	// IsOpenTrace when true, it will record trace info.
	IsOpenTrace bool
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
type defaultLogFormatter struct {
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
	// DecodeURL is a path the client requests.
	DecodeURL string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]any
}

func (f *defaultLogFormatter) BeforeResponse(_ *gin.Context) {}

func (f *defaultLogFormatter) AfterResponse(c *gin.Context) {
	t := trace.GetTrace(c)

	f.StartTime = t.RequestAt
	f.EndTime = t.ResponseAt
	f.Latency = t.Latency
	f.StatusCode = t.Response.HttpCode
	f.ClientIP = t.Request.ClientIP
	f.Method = t.Request.Method
	f.DecodeURL = t.Request.DecodeURL
	f.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
	f.BodySize = c.Writer.Size()
	f.Keys = c.Keys
}

func (f *defaultLogFormatter) Format(c *gin.Context) string {
	if f.Latency > time.Minute {
		f.Latency = f.Latency.Truncate(time.Second)
	}

	return fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %s | %-7s %#v\n%s",
		f.StartTime.Format(time.DateTime),
		f.StatusCode,
		f.Latency,
		f.ClientIP,
		GetRequestIDFromContext(c),
		f.Method,
		f.DecodeURL,
		f.ErrorMessage,
	)
}

func (f *defaultLogFormatter) i() {}

// defaultLogFormatter is the default log format function Logger middleware uses.
//var defaultLogFormatter = func(param LogBasicParams) string {
//}

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

	notlogged := conf.SkipPaths

	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	filterFunc := conf.Filter

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Read the Body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		t := trace.New(GetRequestIDFromContext(c))
		if t == nil {
			return
		}
		c.Set(trace.Header, t)

		t.WithRequest(c)

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
		if _, ok := skip[path]; ok || filterFunc(c) {
			return
		}

		t.WithResponse(c)
		formatter.AfterResponse(c)

		logger.GetLogger().
			WithField("method", t.Request.Method).
			WithField("path", t.Request.DecodeURL).
			WithField("client_ip", t.Request.ClientIP).
			WithField("http_code", t.Response.HttpCode).
			WithField("business_code", t.Response.BusinessCode).
			WithField("business_code_msg", t.Response.BusinessCodeMsg).
			WithField("http_code_msg", t.Response.HttpCodeMsg).
			WithField("trace_id", t.Identifier).
			WithField("request_at", t.RequestAt.Format(time.DateTime)).
			WithField("response_at", t.ResponseAt.Format(time.DateTime)).
			WithField("costs", t.Latency.Microseconds()).
			WithField("sql", t.SQLs).
			Infof(formatter.Format(c))
	}
}
