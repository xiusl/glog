package glog

import (
	"errors"
	"strings"
)

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
