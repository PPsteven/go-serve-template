package context

import (
	"github.com/gin-gonic/gin"
	"go-server-template/pkg/middleware"
)

const (
	_RequestID = middleware.ContextRequestIDKey
)

func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(_RequestID); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}

	return ""
}
