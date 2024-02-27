package context

import (
	"github.com/gin-gonic/gin"
)

const (
	_RequestID = "request_id"
)

func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(_RequestID); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}

	return ""
}
