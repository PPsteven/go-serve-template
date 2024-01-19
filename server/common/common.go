package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func resp(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Resp[interface{}]{
		Code: code,
		Message: msg,
		Data: data,
	})
}

func SuccessResp(c *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		resp(c, SUCCESS, "success", nil)
		return
	}
	resp(c, SUCCESS, "success", data[0])
}

func SuccessWithMessage(c *gin.Context, msg string, data ...interface{}) {
	resp(c, SUCCESS, msg, data)
}

func ErrorResp(c *gin.Context, msg string) {
	resp(c, ERROR, msg, nil)
}

func ErrorRespWithCode(c *gin.Context, msg string, code int) {
	resp(c, code, msg, nil)
}
