package glog

import (
	"fmt"
	"time"
)

func (l *GLogger) Log(level Level, format string, arg ...interface{}) {
	if !l.enable(level) {
		return
	}
	msg := fmt.Sprintf(format, arg...)
	timeStr := time.Now().Format(timeFromatStr)
	filePath, fileName, lineNo := getInfo(3)
	fmt.Printf("[%s] [INFO] [%s:%s:%d] %s\n", timeStr, fileName, filePath, lineNo, msg)
}

func (l *GLogger) Debug(format string, arg ...interface{}) {
	l.Log(DEBUG, format, arg...)
}

func (l *GLogger) Info(format string, arg ...interface{}) {
	l.Log(INFO, format, arg...)
}

func (l *GLogger) Warning(format string, arg ...interface{}) {
	l.Log(WARNING, format, arg...)
}
