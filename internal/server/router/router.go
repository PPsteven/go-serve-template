package router

import (
	"github.com/gin-gonic/gin"
	"go-server-template/internal/server/handlers"
)

func Load(e *gin.Engine, middlewares ...gin.HandlerFunc) {
	{
		e.Use(middlewares...)
		api := e.Group("/api")
		{
			api.GET("/user/:id", handlers.GetUserByID)
		}
	}
}
