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

type Meta map[string]any

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

type LogContext struct {
	File string
	Line int
}

var contextFieldsCount int = reflect.TypeFor[LogContext]().NumField()

func Log(level LogLevel, msg string, meta ...Meta) {
	var metaMap Meta
	if len(meta) > 0 {
		metaMap = meta[0]
	}

	context := LogContext{}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		paths := strings.Split(file, "/")
		context.File = paths[len(paths)-2] + "/" + paths[len(paths)-1]
		context.Line = line
	}

	var length int
	if metaMap != nil {
		length = len(metaMap) * 2
	}

	length += contextFieldsCount * 2

	var logsData = make([]any, 0, length)
	for k, v := range metaMap {
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
