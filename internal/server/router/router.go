package router

import (
	"github.com/gin-gonic/gin"
	"go-server-template/internal/server/handlers"
)

func Init(e *gin.Engine) {
	{
		e.Use()
		api := e.Group("/api")
		{
			api.GET("/user/:id", handlers.GetUserByID)
		}
	}
}
