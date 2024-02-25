package bootstrap

import (
	"go-server-template/internal/db"
	"go-server-template/internal/service"
)

func Init() {
	InitLog()
	InitDB()

	service.Init(db.GetDB())
}
