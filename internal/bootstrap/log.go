package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"go-server-template/internal/conf"
	"io"
	stdLog "log"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
)

func InitLog() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true, // 允许环境变量覆盖彩色日志输出的设置。
		FullTimestamp:             true,
		TimestampFormat:           time.DateTime,
	})

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
		if conf.Conf.Env == conf.Dev {
			w = io.MultiWriter(os.Stdout, w)
		}
		log.SetOutput(w)
	}
	stdLog.SetOutput(log.StandardLogger().Out)

	log.Infof("init logrus...")
}
