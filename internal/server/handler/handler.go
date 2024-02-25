package handler

import (
	"go-server-template/internal/server/handler/api/user"
	"go-server-template/internal/service"
)

func User() user.Handler {
	return user.New(service.Get())
}
