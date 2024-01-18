package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"go-server-template/internal/conf"
	"io"
	stdLog "log"

	"github.com/natefinch/lumberjack"
)

func InitLog() {
	if conf.Conf.Env == conf.Dev {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
		log.SetReportCaller(false)
	}

	logFileConfig := conf.Conf.Logger.LogFile
	if logFileConfig.Enable {
		var w io.Writer = &lumberjack.Logger{
			Filename:   logFileConfig.Name,
			MaxSize:    logFileConfig.MaxSize,
			MaxAge:     logFileConfig.MaxAge,
			MaxBackups: logFileConfig.MaxBackups,
			Compress:   logFileConfig.Compress,
		}
		log.SetOutput(w)
	}
	stdLog.SetOutput(log.StandardLogger().Out)

	log.Infof("init logrus...")
}
