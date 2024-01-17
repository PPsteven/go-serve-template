package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-server-template/internal/conf"
	"go-server-template/internal/model"
	"gorm.io/gorm"
)

var dB *gorm.DB

func Init(d *gorm.DB) {
	dB = d
	err := AutoMigrate(new(model.User))
	if err != nil {
		log.Fatalf("failed migrate database: %s", err.Error())
	}
}

func AutoMigrate(dist ...interface{}) error {
	var err error
	if conf.Conf.Database.Type == "mysql" {
		// TODO ...
	} else {
		err = dB.AutoMigrate(dist...)
	}
	fmt.Println(123)
	return err
}