package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

const (
	// DefaultLevel the default log level
	DefaultLevel = DebugLevel

	// DefaultTimeLayout the default time layout;
	DefaultTimeLayout = time.RFC3339
)

const (
	Logrus = "logrus"
	Zap    = "zap"
)

type Fields map[string]interface{}

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(kvs Fields) Logger

	Sync() error

	Writer() io.Writer

	i()
}

// Option custom setup config
type Option func(*option)

type option struct {
	level            Level
	file             io.Writer
	timeLayout       string
	disableConsole   bool
	disableCaller    bool
	enableStackTrace bool
	encoding         string
}

// WithDebugLevel only greater than 'level' will output
func WithDebugLevel() Option {
	return func(opt *option) {
		opt.level = DebugLevel
	}
}

// WithInfoLevel only greater than 'level' will output
func WithInfoLevel() Option {
	return func(opt *option) {
		opt.level = InfoLevel
	}
}

// WithWarnLevel only greater than 'level' will output
func WithWarnLevel() Option {
	return func(opt *option) {
		opt.level = WarnLevel
	}
}

// WithErrorLevel only greater than 'level' will output
func WithErrorLevel() Option {
	return func(opt *option) {
		opt.level = ErrorLevel
	}
}

// WithFatalLevel only greater than 'level' will output
func WithFatalLevel() Option {
	return func(opt *option) {
		opt.level = FatalLevel
	}
}

func WithDisableCaller() Option {
	return func(opt *option) {
		opt.disableCaller = true
	}
}

func WithEncodingJson() Option {
	return func(opt *option) {
		opt.encoding = "json"
	}
}

// WithFileP write log to some file
func WithFileP(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0766)
	if err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = zapcore.Lock(f)
	}
}

// WithFileRotationP write log to some file with rotation
func WithFileRotationP(file string, maxSize, maxBackups, maxAge int, localTime, compress bool) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = &lumberjack.Logger{ // concurrent-safed
			Filename:   file,       // 文件路径
			MaxSize:    maxSize,    // 单个文件最大尺寸，默认单位 M
			MaxBackups: maxBackups, // 最多保留 300 个备份
			MaxAge:     maxAge,     // 最大时间，默认单位 day
			LocalTime:  localTime,  // 使用本地时间
			Compress:   compress,   // 是否压缩 disabled by default
		}
	}
}

// WithTimeLayout custom time format
func WithTimeLayout(timeLayout string) Option {
	return func(opt *option) {
		opt.timeLayout = timeLayout
	}
}

// WithDisableConsole WithEnableConsole write log to os.Stdout or os.Stderr
func WithDisableConsole() Option {
	return func(opt *option) {
		opt.disableConsole = true
	}
}

var globalLog Logger

func GetLogger() Logger {
	return globalLog
}

// Init init logger
func Init(logType string, options ...Option) Logger {
	var log Logger
	if logType == "zap" {
		log, _ = NewZapLogger(options...)
	} else {
		log, _ = NewZapLogger(options...)
	}
	globalLog = log
	return log
}
