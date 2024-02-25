package response

import (
	"github.com/gin-gonic/gin"
	"go-server-template/internal/server/errcode"
	"net/http"
)

type Response struct {
	Code    int
	Message string
	Data    interface{}
	Detail  []string
}

func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}

	c.JSON(http.StatusOK, data)
}

func SuccessWithCode(c *gin.Context, data interface{}, code int) {
	if data == nil {
		data = gin.H{}
	}

	if code < http.StatusOK || code > http.StatusIMUsed {
		code = http.StatusOK
	}

	c.JSON(code, data)
}

func Error(c *gin.Context, err errcode.SvrError) {
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	response := &Response{
		Code:    err.Code(),
		Message: err.Message(),
		Data:    gin.H{},
		Detail:  err.Detail(),
	}

	c.JSON(err.HttpCode(), response)
}
