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

type Logger struct {
	LogLevel string `json:"log_level"`
}

var Conf *Config

func InitDefaultConfig() *Config {
	dbFile := filepath.Join("data", "data.db")
	return &Config{
		Database: Database{
			Type: "sqlite3",
			Port: 0,
			File: dbFile,
		},
		Logger: Logger{
			LogLevel: "debug",
		},
		Env: Dev,
	}
}
