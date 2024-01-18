package db

import (
	"gorm.io/gorm"
)

var dB *gorm.DB

func InitDB(d *gorm.DB) {
	dB = d
}

func GetDB() *gorm.DB {
	return dB
}