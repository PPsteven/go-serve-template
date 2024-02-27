package response

import (
	"github.com/gin-gonic/gin"
	"go-server-template/internal/server/errcode"
	"go-server-template/pkg/context"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Detail  []string    `json:"detail"`
	TraceID string      `json:"trace_id"`
}

func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}

	response := &Response{
		Code:    0,
		Message: "success",
		Data:    data,
		Detail:  []string{},
		TraceID: context.GetRequestID(c),
	}

	c.JSON(http.StatusOK, response)
}

func SuccessWithHttpCode(c *gin.Context, data interface{}, httpCode int) {
	if data == nil {
		data = gin.H{}
	}

	if httpCode < http.StatusOK || httpCode > http.StatusIMUsed {
		httpCode = http.StatusOK
	}

	response := &Response{
		Code:    0,
		Message: "success",
		Data:    data,
		Detail:  []string{},
		TraceID: context.GetRequestID(c),
	}

	c.JSON(httpCode, response)
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
		TraceID: context.GetRequestID(c),
	}

	c.JSON(err.HttpCode(), response)
}
