package db

import (
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(d *gorm.DB) {
	db = d
}

func GetDB() *gorm.DB {
	return db
}