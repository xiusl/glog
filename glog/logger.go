package glog

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

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
