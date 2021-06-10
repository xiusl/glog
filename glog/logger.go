package glog

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

type GLogger interface {
	Debug(format string, arg ...interface{})
	Info(format string, arg ...interface{})
	Trace(format string, arg ...interface{})
	Warning(format string, arg ...interface{})
	Error(format string, arg ...interface{})
	Fatal(format string, arg ...interface{})
}

type GLoggerConfig struct {
	Type     string
	Level    string
	FilePath string
	FileName string
}

func NewLogger(config *GLoggerConfig) GLogger {
	switch config.Type {
	case "file":
		return NewFileLogger(config.Level, config.FilePath, config.FileName)
	case "console":
		return NewConsoleLogger(config.Level)
	default:
		panic(errors.New("Logger type only support `file` `console`"))
	}
}

const timeFromatStr = "2006-01-02 15:04:05"

func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime caller failed")
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	funcName = strings.Split(funcName, ".")[1]
	return
}
