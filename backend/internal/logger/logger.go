package logger

import "log"

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	WARN  LogLevel = "WARN"
)

const (
	infoPrefix  = "[INFO] "
	warnPrefix  = "[WARN] "
	errorPrefix = "[ERROR] "
	debugPrefix = "[DEBUG] "
)

var logger = log.New(log.Writer(), "", log.LstdFlags|log.Lshortfile)

func Log(level LogLevel, fmt string, args ...any) {
	switch level {
	case DEBUG:
		logger.Printf(debugPrefix+fmt, args...)
	case INFO:
		logger.Printf(infoPrefix+fmt, args...)
	case WARN:
		logger.Printf(warnPrefix+fmt, args...)
	case ERROR:
		logger.Printf(errorPrefix+fmt, args...)
	}
}
