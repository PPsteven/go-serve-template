package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-server-template/internal/conf"
	"go-server-template/internal/server/errcode"
	"go-server-template/internal/server/response"
	"go-server-template/pkg/app"
	"go-server-template/pkg/context"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Parse the token.
		header := c.Request.Header.Get("Authorization")

		if len(header) == 0 {
			response.Error(c, errcode.ErrInvalidAuthorization)
			c.Abort()
			return
		}

		var t string
		_, err := fmt.Sscanf(header, "Bearer %s", &t)
		if err != nil {
			response.Error(c, errcode.ErrInvalidAuthorization)
			c.Abort()
			return
		}

		// Parse the json web token
		ctx, err := app.Parse(t, conf.Conf.JWT.Secret)
		if err != nil {
			response.Error(c, errcode.ErrInvalidAuthorization)
			c.Abort()
			return
		}

		context.SetUserID(c, ctx.UserID)

		c.Next()
	}
}
