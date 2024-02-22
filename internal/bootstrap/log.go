package bootstrap

import (
	"go-server-template/internal/conf"
	"go-server-template/pkg/logger"
)

func InitLog() {
	opts := []logger.Option{
		logger.WithEncodingJson(), // json format
	}

	if conf.Conf.Env == conf.Production {
		opts = append(opts, logger.WithInfoLevel())
		opts = append(opts, logger.WithDisableCaller())
	} else {
		opts = append(opts, logger.WithDebugLevel())
	}

	logFileConfig := conf.Conf.Logger.LogFile
	if logFileConfig.Enable {
		opts = append(opts, logger.WithFileRotationP(
			logFileConfig.Name,
			logFileConfig.MaxSize,
			logFileConfig.MaxBackups,
			logFileConfig.MaxAge,
			logFileConfig.LocalTime,
			logFileConfig.Compress,
		))
	}

	log, _ := logger.NewZapLogger(opts...)
	logger.InitGolbalLogger(log)

	logger.GetLogger().Infof("log init success.")
}
