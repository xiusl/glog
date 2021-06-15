package logging

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

type Logginger interface {
	Debug(format string, arg ...interface{})
	Info(format string, arg ...interface{})
	Trace(format string, arg ...interface{})
	Warning(format string, arg ...interface{})
	Error(format string, arg ...interface{})
	Fatal(format string, arg ...interface{})
}

type LoggingConfig struct {
	Type     string
	Level    string
	FilePath string
	FileName string
}

func NewLogger(config *LoggingConfig) Logginger {
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

// Level 日志的级别
type Level uint16

const (
	UNKOWN Level = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

func (l Level) ToString() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKOWN"
	}
}
func parseLevel(s string) (Level, error) {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("Invalid log level")
		return UNKOWN, err
	}
}
