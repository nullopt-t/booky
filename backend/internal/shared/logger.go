package shared

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

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

var logger = log.New(log.Writer(), "", log.LstdFlags)

func Log(level LogLevel, format string, args ...any) {
	var sfile string = "nil"
	pointer, file, line, ok := runtime.Caller(1)
	if ok {
		fn := file[strings.LastIndex(file, "/"):]
		sfile = fmt.Sprintf("0x%x %s:%d", pointer, fn, line)
	}

	switch level {
	case DEBUG:
		logger.Printf(debugPrefix+sfile+format, args...)
	case INFO:
		logger.Printf(infoPrefix+sfile+format, args...)
	case WARN:
		logger.Printf(warnPrefix+sfile+format, args...)
	case ERROR:
		logger.Printf(errorPrefix+sfile+format, args...)
	}
}
