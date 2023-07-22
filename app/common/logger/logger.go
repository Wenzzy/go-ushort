package logger

import (
	"encoding/json"
	"path"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wenzzyx/go-ushort/app/config"
)

var logger = logrus.New()

func loggerPrettyfier(frame *runtime.Frame) (function string, file string) {
	fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
	return "", fileName
}

func InitLogger() {
	cfg := config.GetCfg()
	logger.Level = logrus.InfoLevel
	if cfg.Server.IsDebug {
		logger.SetFormatter(&logrus.TextFormatter{CallerPrettyfier: loggerPrettyfier})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{CallerPrettyfier: loggerPrettyfier})
	}

	logger.SetReportCaller(true)
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

func JsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006-01-02T15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()
			log["referrer"] = params.Request.Referer()
			log["request_id"] = params.Request.Header.Get("Request-Id")

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}
