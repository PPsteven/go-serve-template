package context

import (
	"github.com/gin-gonic/gin"
)

const (
	_RequestID = "request_id"
	_Alias     = "_alias_"
	_UserID    = "_user_id_"
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

func SetUserID(c *gin.Context, userID uint64) {
	c.Set(_UserID, userID)
}

func GetUserID(c *gin.Context) uint64 {
	if v, ok := c.Get(_UserID); ok {
		if userID, ok := v.(uint64); ok {
			return userID
		}
	}
	return 0
}