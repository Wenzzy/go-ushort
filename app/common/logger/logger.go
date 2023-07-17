package logger

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var logger = logrus.New()

func init() {
	logger.Level = logrus.InfoLevel
	logger.Formatter = &formatter{}

	logger.SetReportCaller(true)
}

func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

type Fields logrus.Fields

func Debugf(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Fatalf(format, args...)
	}
}

type formatter struct {
	prefix string
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer

	sb.WriteString(strings.ToUpper(entry.Level.String()))
	sb.WriteString(" ")
	sb.WriteString(entry.Time.Format(time.RFC3339))
	sb.WriteString(" ")
	sb.WriteString(f.prefix)
	sb.WriteString(entry.Message)

	return sb.Bytes(), nil
}
