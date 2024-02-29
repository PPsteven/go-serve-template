package router

import (
	"github.com/gin-gonic/gin"
	"go-server-template/internal/server/handlers"
	"go-server-template/pkg/middleware"
)

func Load(e *gin.Engine, middlewares ...gin.HandlerFunc) {
	{
		e.Use(middlewares...)
		// api := e.Group("/api", middleware_internal.Auth())
		api := e.Group("/api")
		{
			api.GET("/user/:id", middleware.Alias("/user/:id"), handlers.User().GetUser)
		}
	}
}
