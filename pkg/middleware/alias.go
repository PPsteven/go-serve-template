package middleware

import (
	"github.com/gin-gonic/gin"
	"go-server-template/pkg/context"
)

// Alias set alias for request url
func Alias(alias string) gin.HandlerFunc {
	return func(c *gin.Context) {
		context.SetAlias(c, alias)

		c.Next()
	}
}
