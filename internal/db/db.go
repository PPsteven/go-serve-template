package db

import (
	log "github.com/sirupsen/logrus"
	"go-server-template/internal/conf"
	"gorm.io/gorm"
)

var dB *gorm.DB

func Init(d *gorm.DB) {
	dB = d
	err := AutoMigrate()
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
	return err
}