package bootstrap

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-server-template/internal/conf"
	"go-server-template/internal/db"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	stdlog "log"
	"strings"
	"time"
)

func InitDB() {
	var (
		dB       *gorm.DB
		err      error
		logLevel logger.LogLevel
	)

	config := conf.Conf

	if config.Env == conf.Dev {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	gormLogger := logger.New(
		stdlog.New(log.StandardLogger().Out, "\r\n", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	gormConfig := &gorm.Config{
		Logger: gormLogger,
	}

	database := config.Database
	switch database.Type {
	case "sqlite3":
		if !(strings.HasSuffix(database.File, ".db") && len(database.File) > 3) {
			log.Fatalf("db name error.")
		}
		dB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental",
			database.File)), gormConfig)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
			database.User, database.Password, database.Host, database.Port, database.Name, database.SSLMode)
		dB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
			database.Host, database.User, database.Password, database.Name, database.Port, database.SSLMode)
		dB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	default:
		log.Fatalf("not supported database type: %s", database.Type)
	}
	if err != nil {
		log.Fatalf("failed to connect database: %s", err.Error())
	}

	db.Init(dB)
}
