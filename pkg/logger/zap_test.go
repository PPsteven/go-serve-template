package logger

import (
	"errors"
	"testing"
	"time"
)

func TestZapLogger(t *testing.T) {
	logger, err := NewZapLogger()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	err = errors.New("err message")

	t.Run("console", func(t *testing.T) {
		logger.WithField("para1", "value1").WithField("para2", "value2").Error(err)
		logger.WithFields(Fields{"para1": "value1", "para2": "value2"}).Info(err)
	})

	t.Run("json", func(t *testing.T) {
		logger, _ = NewZapLogger(WithEncodingJson())
		defer logger.Sync()

		logger.WithField("para1", "value1").WithField("para2", "value2").Error(err)
		logger.WithFields(Fields{"para1": "value1", "para2": "value2"}).Info(err)
	})

	t.Run("time layout", func(t *testing.T) {
		logger, _ = NewZapLogger(WithTimeLayout(time.DateTime))
		defer logger.Sync()

		logger.WithField("para1", "value1").WithField("para2", "value2").Error(err)
		logger.WithFields(Fields{"para1": "value1", "para2": "value2"}).Info(err)
	})

	t.Run("level", func(t *testing.T) {
		logger, _ = NewZapLogger(WithErrorLevel())
		defer logger.Sync()

		logger.WithField("para1", "value1").WithField("para2", "value2").Error(err)
		logger.WithFields(Fields{"para1": "value1", "para2": "value2"}).Info(err)
	})
}

func TestZapLoggerWithFile(t *testing.T) {
	logger, err := NewZapLogger(WithFileP("test.log"), WithEncodingJson())
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	err = errors.New("err message")

	logger.WithField("para1", "value1").WithField("para2", "value2").
		Info(err)
}
