package main

import (
	"time"

	"github.com/xiusl/glog/glog"
)

func main() {
	go alog()
	blog()
}

func alog() {
	logConf := &glog.GLoggerConfig{
		Type:     "file",
		FilePath: "./logs",
		FileName: "a.log",
		Level:    "debug",
	}

	log := glog.NewLogger(logConf)

	for {
		log.Debug("aaaa -> this is a debug log msg")
		log.Warning("aaaa -> this is a warning log msg")

		time.Sleep(10 * time.Millisecond)
	}
}

func blog() {
	logConf := &glog.GLoggerConfig{
		Type:     "file",
		FilePath: "./logs",
		FileName: "b.log",
		Level:    "debug",
	}

	log := glog.NewLogger(logConf)

	for {
		log.Debug("bbbb -> this is a debug log msg")
		log.Warning("bbbb -> this is a warning log msg")

		time.Sleep(500 * time.Millisecond)
	}
}
