package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var _ Logger = (*logrusEntry)(nil)

type logrusEntry struct {
	entry *logrus.Entry
}

func (e *logrusEntry) i() {}

func (e *logrusEntry) Debug(args ...interface{}) {
	e.entry.Debug(args...)
}

func (e *logrusEntry) Debugf(format string, args ...interface{}) {
	e.entry.Debugf(format, args...)
}

func (e *logrusEntry) Info(args ...interface{}) {
	e.entry.Info(args...)
}

func (e *logrusEntry) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

func (e *logrusEntry) Warn(args ...interface{}) {
	e.entry.Warn(args...)
}

func (e *logrusEntry) Warnf(format string, args ...interface{}) {
	e.entry.Warnf(format, args...)
}

func (e *logrusEntry) Error(args ...interface{}) {
	e.entry.Error(args...)
}

func (e *logrusEntry) Errorf(format string, args ...interface{}) {
	e.entry.Errorf(format, args...)
}

func (e *logrusEntry) WithField(key string, value interface{}) Logger {
	clone := e.clone()
	clone.entry = clone.entry.WithField(key, value)
	return e
}

func (e *logrusEntry) WithFields(kvs Fields) Logger {
	var data = make(logrus.Fields, len(kvs))
	for k, v := range kvs {
		data[k] = v
	}
	clone := e.clone()
	clone.entry = clone.entry.WithFields(data)
	return e
}

func (e *logrusEntry) Sync() error {
	return nil
}

func (e *logrusEntry) Writer() io.Writer {
	return e.entry.Writer()
}

func (e *logrusEntry) clone() *logrusEntry {
	c := *e
	return &c
}

func NewLogrusLogger(opts ...Option) (Logger, error) {
	opt := &option{level: DefaultLevel, encoding: EncodingConsole}
	for _, f := range opts {
		f(opt)
	}

	timeLayout := DefaultTimeLayout
	if opt.timeLayout != "" {
		timeLayout = opt.timeLayout
	}

	var formatter logrus.Formatter

	if opt.encoding == EncodingConsole {
		formatter = &logrus.TextFormatter{
			TimestampFormat: timeLayout,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "time",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "msg",
				logrus.FieldKeyFile:  "caller",
			},
		}
	} else {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: timeLayout,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "time",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "msg",
				logrus.FieldKeyFile:  "caller",
			},
		}
	}

	logger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(logrus.LevelHooks),
		Level:        logrusLevel(opt.level),
		ExitFunc:     os.Exit,
		ReportCaller: !opt.disableCaller,
	}

	if opt.file != nil {
		logger.Out = io.MultiWriter(logger.Out, opt.file)
	}

	entry := logger.WithContext(context.TODO())

	return &logrusEntry{entry: entry}, nil
}

func logrusLevel(level Level) logrus.Level {
	switch level {
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
