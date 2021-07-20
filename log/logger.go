package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Log interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

var (
	defaultOut       = os.Stdout
	defaultLogLevel  = "debug"
	defaultFormatter = &logrus.TextFormatter{}
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger() (logger *Logger) {
	logger = &Logger{
		log: logrus.New(),
	}
	logger.log.SetOutput(defaultOut)
	logger.log.Out = defaultOut
	logger.log.Formatter = &logrus.JSONFormatter{}
	return
}

func (l *Logger) SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		l.log.SetLevel(logrus.DebugLevel)
	case "info":
		l.log.SetLevel(logrus.InfoLevel)
	case "warn", "warning":
		l.log.SetLevel(logrus.WarnLevel)
	case "error":
		l.log.SetLevel(logrus.ErrorLevel)
	default:
		panic("log level must be: debug、info、warn、error")
	}
}

func (l *Logger) SetOutFile(filePath string) {
	src, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic("create log fill error : " + err.Error())
	}
	l.log.Out = src
}

func (l *Logger) WithFields(fields map[string]interface{}) *Entry {
	return NewEntity(l.log, fields)
}

func (l *Logger) Debug(v ...interface{}) {
	l.log.Debug(v)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.log.Debugf(format, v)
}

func (l *Logger) Info(v ...interface{}) {
	l.log.Info(v)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.log.Infof(format, v)
}

func (l *Logger) Warn(v ...interface{}) {
	l.log.Warn(v)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.log.Warnf(format, v)
}

func (l *Logger) Error(v ...interface{}) {
	l.log.Error(v)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.log.Errorf(format, v)
}
