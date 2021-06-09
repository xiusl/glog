package glog

import (
	"fmt"
	"time"
)

type ConsoleLogger struct {
	Level Level
}

func NewLogger(levelStr string) *ConsoleLogger {
	level, err := parseLevel(levelStr)
	if err != nil {
		panic(err)
	}

	return &ConsoleLogger{
		Level: level,
	}
}

func (l *ConsoleLogger) enable(level Level) bool {
	return l.Level <= level
}

func (l *ConsoleLogger) Log(level Level, format string, arg ...interface{}) {
	if !l.enable(level) {
		return
	}
	msg := fmt.Sprintf(format, arg...)
	timeStr := time.Now().Format(timeFromatStr)
	filePath, fileName, lineNo := getInfo(3)
	fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", timeStr, level.ToString(), fileName, filePath, lineNo, msg)
}

func (l *ConsoleLogger) Debug(format string, arg ...interface{}) {
	l.Log(DEBUG, format, arg...)
}

func (l *ConsoleLogger) Info(format string, arg ...interface{}) {
	l.Log(INFO, format, arg...)
}

func (l *ConsoleLogger) Warning(format string, arg ...interface{}) {
	l.Log(WARNING, format, arg...)
}

func (l *ConsoleLogger) Error(format string, arg ...interface{}) {
	l.Log(ERROR, format, arg...)
}

func (l *ConsoleLogger) Fatal(format string, arg ...interface{}) {
	l.Log(FATAL, format, arg...)
}
