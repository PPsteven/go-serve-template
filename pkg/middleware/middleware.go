package middleware

import "github.com/gin-gonic/gin"

// Middlewares global middleware
var Middlewares = defaultMiddlewares()

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":   gin.Recovery(),
		"secure":     Secure,
		"options":    Options,
		"nocache":    NoCache,
		"request_id": RequestID(),
		"logger":     Logger(),
		//"trace":      Tracing,
		//"cors":       Cors(),
	}
}
