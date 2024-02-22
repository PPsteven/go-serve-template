package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

const (
	// EncodingConsole console输出
	EncodingConsole = "console"
	// EncodingJson json输出
	EncodingJson = "json"
)

// For mapping config logger to app logger levels
var loggerLevelMap = map[Level]zapcore.Level{
	DebugLevel: zap.DebugLevel,
	InfoLevel:  zap.InfoLevel,
	WarnLevel:  zap.WarnLevel,
	ErrorLevel: zap.ErrorLevel,
	FatalLevel: zap.FatalLevel,
}

var _ Logger = (*zapEntry)(nil)

type zapEntry struct {
	entry *zap.Logger
}

func (e *zapEntry) i() {}

func (e *zapEntry) Debug(args ...interface{}) {
	e.entry.Debug(fmt.Sprint(args...))
}

func (e *zapEntry) Debugf(format string, args ...interface{}) {
	e.entry.Debug(fmt.Sprintf(format, args...))
}

func (e *zapEntry) Info(args ...interface{}) {
	e.entry.Info(fmt.Sprint(args...))
}

func (e *zapEntry) Infof(format string, args ...interface{}) {
	e.entry.Info(fmt.Sprintf(format, args...))
}

func (e *zapEntry) Warn(args ...interface{}) {
	e.entry.Warn(fmt.Sprint(args...))
}

func (e *zapEntry) Warnf(format string, args ...interface{}) {
	e.entry.Warn(fmt.Sprintf(format, args...))
}

func (e *zapEntry) Error(args ...interface{}) {
	e.entry.Error(fmt.Sprint(args...))
}

func (e *zapEntry) Errorf(format string, args ...interface{}) {
	e.entry.Error(fmt.Sprintf(format, args...))
}

func (e *zapEntry) WithField(key string, value interface{}) Logger {
	clone := e.clone()
	clone.entry = clone.entry.With(zap.Any(key, value))
	return clone
}

func (e *zapEntry) WithFields(kvs Fields) Logger {
	var data = make([]zap.Field, 0, len(kvs))
	for k, v := range kvs {
		data = append(data, zap.Any(k, v))
	}
	clone := e.clone()
	clone.entry = clone.entry.With(data...)
	return clone
}

func (e *zapEntry) Sync() error {
	return e.entry.Sync()
}

func (e *zapEntry) Writer() io.Writer {
	return zap.NewStdLog(e.entry).Writer()
}

func (e *zapEntry) clone() *zapEntry {
	c := *e
	return &c
}

// NewZapLogger new zap logger
func NewZapLogger(opts ...Option) (Logger, error) {
	opt := &option{level: DefaultLevel, encoding: EncodingConsole}
	for _, f := range opts {
		f(opt)
	}
	curLevel := zapLevel(opt.level)

	var encoderCfg zapcore.EncoderConfig
	var enc zapcore.Encoder

	encoderCfg = newEncoderConfig(opt)

	if opt.encoding == EncodingConsole {
		enc = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		enc = zapcore.NewJSONEncoder(encoderCfg)
	}

	stdout := zapcore.Lock(os.Stdout)
	stderr := zapcore.Lock(os.Stderr)

	var cores []zapcore.Core
	var options []zap.Option

	options = append(options, zap.ErrorOutput(stderr))

	// add stacktrace
	if opt.enableStackTrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// add caller
	if !opt.disableCaller {
		options = append(options, zap.AddCaller())
	}

	// output to console
	if !opt.disableConsole {
		cores = append(cores, zapcore.NewCore(
			enc,
			zapcore.NewMultiWriteSyncer(stdout),
			curLevel,
		))
	}

	// output to file
	if opt.file != nil {
		cores = append(cores, zapcore.NewCore(
			enc,
			zapcore.AddSync(opt.file),
			curLevel,
		))
	}

	combinedCore := zapcore.NewTee(cores...)

	logger := zap.New(combinedCore, options...)

	return &zapEntry{logger}, nil
}

func zapLevel(level Level) zapcore.Level {
	lvl, exists := loggerLevelMap[level]
	if !exists {
		return zapcore.InfoLevel
	}

	return lvl
}

// similar to zap.NewProductionEncoderConfig
// newEncoderConfig returns a customized encoder config
func newEncoderConfig(opt *option) zapcore.EncoderConfig {
	timeLayout := DefaultTimeLayout
	if opt.timeLayout != "" {
		timeLayout = opt.timeLayout
	}

	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger", // used by logger.Named(key); optional; useless
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeLayout), // time format by layout
		EncodeDuration: zapcore.MillisDurationEncoder,           // integer number of milliseconds elapsed
		EncodeCaller:   zapcore.ShortCallerEncoder,              // package/file:line format
	}
}
