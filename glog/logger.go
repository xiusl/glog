package glog

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

const timeFromatStr = "2006-01-02 15:04:05"

type GLogger struct {
	Level Level
}

func NewLogger(levelStr string) *GLogger {
	level, err := parseLevel(levelStr)
	if err != nil {
		panic(err)
	}

	return &GLogger{
		Level: level,
	}
}

func (l *GLogger) enable(level Level) bool {
	return l.Level <= level
}

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
