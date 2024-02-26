// Copy From gin-gonic/gin/logger.go
// Aim to support add more custom params in log
// Alternative: https://github.com/gin-contrib/zap

package middleware

import (
	"bytes"
	"go-server-template/pkg/logger"
	"go-server-template/pkg/trace"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string

	Filter func(ctx *gin.Context) bool

	// IsOpenTrace when true, it will record trace info.
	IsOpenTrace bool
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

		// Continue.
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; ok || filterFunc(c) {
			return
		}

		t.WithResponse(c, blw.body)

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
			Info("trace info")
	}
}
