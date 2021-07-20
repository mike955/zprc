package log

func LoggerHelper(logger *Logger, fields map[string]interface{}) *Entry {
	return NewEntity(logger.log, fields)
}

func Helper(logger *Entry, fields map[string]interface{}) *Entry {
	return logger.WithField(fields)
}
