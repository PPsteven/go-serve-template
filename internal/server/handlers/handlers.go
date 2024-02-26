package handlers

import (
	"go-server-template/internal/server/handlers/api/user"
	"go-server-template/internal/service"
)

func User() user.Handler {
	return user.New(service.Get())
}
