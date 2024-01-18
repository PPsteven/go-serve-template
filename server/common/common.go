package common

import "github.com/gin-gonic/gin"

type Resp[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func SuccessResp(c *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, Resp[interface{}]{
			Code: 200,
			Message: "success",
			Data: nil,
		})
		return
	}
	c.JSON(200, Resp[interface{}]{
		Code: 200,
		Message: "success",
		Data: data[0],
	})
}

func ErrorResp(c *gin.Context, err error, code int) {
	c.JSON(200, Resp[interface{}]{
		Code: code,
		Message: err.Error(),
		Data: nil,
	})
}
