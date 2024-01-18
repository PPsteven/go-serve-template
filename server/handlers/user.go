package handlers

import (
	"github.com/gin-gonic/gin"
	"go-server-template/server/common"
	"strconv"
)

func GetUserByID(c *gin.Context) {
	_, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
	}

	common.SuccessResp(c, nil)
}