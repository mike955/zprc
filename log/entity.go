package log

import (
	"github.com/sirupsen/logrus"
)

type Entry struct {
	entry *logrus.Entry
}

func NewEntity(log *logrus.Logger, fileds map[string]interface{}) (entity *Entry) {
	e := log.WithFields(fileds)
	entity = &Entry{
		entry: e,
	}
	return
}

func (e *Entry) Debug(v ...interface{}) {
	e.entry.Debug(v)
}

func (e *Entry) Debugf(format string, v ...interface{}) {
	e.entry.Debugf(format, v)
}

func (e *Entry) Info(v ...interface{}) {
	e.entry.Info(v)
}

func (e *Entry) Infof(format string, v ...interface{}) {
	e.entry.Infof(format, v)
}

func (e *Entry) Warn(v ...interface{}) {
	e.entry.Warn(v)
}

func (e *Entry) Warnf(format string, v ...interface{}) {
	e.entry.Warnf(format, v)
}

func (e *Entry) Error(v ...interface{}) {
	e.entry.Error(v)
}

func (e *Entry) Errorf(format string, v ...interface{}) {
	e.entry.Errorf(format, v)
}

func (e *Entry) WithField(fileds map[string]interface{}) *Entry {
	return &Entry{
		entry: e.entry.WithFields(fileds),
	}
}
