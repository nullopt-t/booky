package shared

import (
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strings"
	"unicode/utf8"
)

type LogLevel int

const (
	DEBUG = iota
	INFO
	ERROR
	WARN
)

type Meta map[string]string

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type LogContext struct {
	File string
	Line int
}

func Log(level LogLevel, msg string, meta Meta) {
	context := LogContext{}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		context.File = file[strings.LastIndex(file, "/"):]
		context.Line = line
	}

	var logsData = make([]any, 0, len(meta)*2+reflect.TypeFor[LogContext]().NumField()*2)
	for k, v := range meta {
		logsData = append(logsData, k, v)
	}

	if utf8.RuneCountInString(context.File) != 0 {
		logsData = append(logsData, "File", context.File)
		logsData = append(logsData, "Line", context.Line)
	}

	switch level {
	case DEBUG:
		logger.Debug(msg, logsData...)
	case INFO:
		logger.Info(msg, logsData...)
	case WARN:
		logger.Warn(msg, logsData...)
	case ERROR:
		logger.Error(msg, logsData...)
	}
}
