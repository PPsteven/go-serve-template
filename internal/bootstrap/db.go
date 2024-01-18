package bootstrap

import (
	"fmt"
	"go-server-template/internal/conf"
	"go-server-template/internal/db"
	"go-server-template/internal/model"
	stdlog "log"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
			SlowThreshold:             200 * time.Millisecond, // 慢查询阈值
			LogLevel:                  logLevel,               // 日志等级
			IgnoreRecordNotFoundError: false,                  // 忽略RecordNotFound错误
			Colorful:                  true,                   // 显示彩色
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
	db.InitDB(dB)
	registerTables()
}

func registerTables() {
	err := AutoMigrate(new(model.User))
	if err != nil {
		log.Fatalf("failed migrate database: %s", err.Error())
	}
	log.Info("register table success")
}

func AutoMigrate(dist ...interface{}) error {
	var err error
	if conf.Conf.Database.Type == "mysql" {
		// TODO ...
	} else {
		err = db.GetDB().AutoMigrate(dist...)
	}
	return err
}
