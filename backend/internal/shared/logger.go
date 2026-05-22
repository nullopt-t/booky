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

var contextFieldsCount int = reflect.TypeFor[LogContext]().NumField()

func Log(level LogLevel, msg string, meta Meta) {
	context := LogContext{}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		context.File = file[strings.LastIndex(file, "/"):]
		context.Line = line
	}

	var length int
	if meta != nil {
		length = len(meta) * 2
	}

	length += contextFieldsCount * 2
	
	var logsData = make([]any, 0, length)
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
