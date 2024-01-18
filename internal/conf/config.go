package conf

import "path/filepath"

type EnvMode string

const (
	Dev        EnvMode = "dev"
	Production         = "production"
)

type Config struct {
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
	Env      EnvMode  `json:"env"`
	Port     int      `json:"port"`
}

type Database struct {
	Type     string `json:"type" env:"DB_TYPE"`
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	User     string `json:"user" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Name     string `json:"name" env:"DB_NAME"`

	File    string `json:"file" env:"DB_PATH"`
	SSLMode string `json:"ssl_mode" env:"DB_SSL_MODE"`
	DSN     string `json:"dsn" env:"DB_DSN"`
}

type LogFile struct {
	Enable     bool   `json:"enable" env:"LOG_ENABLE"`
	Name       string `json:"name" env:"LOG_NAME"`
	MaxSize    int    `json:"max_size" env:"MAX_SIZE"`
	MaxBackups int    `json:"max_backups" env:"MAX_BACKUPS"`
	MaxAge     int    `json:"max_age" env:"MAX_AGE"`
	Compress   bool   `json:"compress" env:"COMPRESS"`
}

type Logger struct {
	LogLevel string  `json:"log_level"`
	LogFile  LogFile `json:"file"`
}

var Conf *Config

func InitDefaultConfig() *Config {
	dbFile := filepath.Join("data", "data.db")
	logPath := filepath.Join("log", "default.log")
	return &Config{
		Database: Database{
			Type: "sqlite3",
			Port: 0,
			File: dbFile,
		},
		Logger: Logger{
			LogLevel: "debug",
			LogFile: LogFile{
				Enable:     true,
				Name:       logPath,
				MaxSize:    10,
				MaxBackups: 5,
				MaxAge:     28,
				Compress:   false,
			},
		},
		Env: Dev,
	}
}
