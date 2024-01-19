package handlers

import (
	"github.com/gin-gonic/gin"
	"go-server-template/internal/db"
	"go-server-template/server/common"
	"strconv"
)

func GetUserByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ErrorResp(c, err.Error())
		return
	}
	user, err := db.GetUserByID(uint(userID))
	if err != nil {
		common.ErrorResp(c, err.Error())
		return
	}

	common.SuccessResp(c, user)
}