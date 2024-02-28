package context

import (
	"github.com/gin-gonic/gin"
)

const (
	_RequestID = "request_id"
	_Alias     = "_alias_"
)

func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(_RequestID); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}

	return ""
}

func SetAlias(c *gin.Context, alias string) {
	c.Set(_Alias, alias)
}

func GetAlias(c *gin.Context) string {
	if v, ok := c.Get(_Alias); ok {
		if alias, ok := v.(string); ok {
			return alias
		}
	}
	return ""
}
