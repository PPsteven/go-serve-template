package logger

import (
	"errors"
	"testing"
	"time"
)

func TestLogrusLogger(t *testing.T) {
	logger, err := NewLogrusLogger()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	err = errors.New("pkg error")

	t.Run("console", func(t *testing.T) {
		logger.WithField("para1", "value1").WithField("para2", "value2").Error(err)
		logger.WithFields(Fields{"para1": "value1", "para2": "value2"}).Info(err)
	})

	t.Run("json", func(t *testing.T) {
		logger, _ = NewLogrusLogger(WithEncodingJson())

		logger.WithField("para1", "value1").WithField("para2", "value2").Error(err)
		logger.WithFields(Fields{"para1": "value1", "para2": "value2"}).Info(err)
	})

	t.Run("time layout", func(t *testing.T) {
		logger, _ = NewLogrusLogger(WithTimeLayout(time.DateTime))

		logger.WithField("para1", "value1").WithField("para2", "value2").Error(err)
		logger.WithFields(Fields{"para1": "value1", "para2": "value2"}).Info(err)
	})

	t.Run("level", func(t *testing.T) {
		logger, _ = NewLogrusLogger(WithErrorLevel())

		logger.WithField("para1", "value1").WithField("para2", "value2").Error(err)
		logger.WithFields(Fields{"para1": "value1", "para2": "value2"}).Info(err)
	})
}

func TestLogrusLoggerWithFile(t *testing.T) {
	// Create a new Logrus logger with a file output
	logger, err := NewLogrusLogger(WithFileP("test.log"), WithEncodingJson())
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	err = errors.New("pkg error")
	logger.WithField("param1", "value1").WithField("param2", "value2").Error(err)
}
