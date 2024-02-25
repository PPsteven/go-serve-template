package user

import (
	"github.com/gin-gonic/gin"
	"go-server-template/internal/server/errcode"
	"go-server-template/internal/server/response"
	"go-server-template/internal/service"
	"strconv"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	GetUser(c *gin.Context)

	i()
}

type handler struct {
	userService service.UserService
}

func New(s service.Service) Handler {
	return &handler{
		userService: s.User(),
	}
}

func (h *handler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParams.WithDetail("invalid user id: %v", err))
		return
	}

	user, err := h.userService.GetUserByID(c, uint(userID))
	if err != nil {
		response.Error(c, errcode.ErrUserNotFound.WithError(err))
		return
	}

	response.Success(c, user)
}

func (h *handler) i() {}
