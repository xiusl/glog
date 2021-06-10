package main

import (
	"fmt"

	"github.com/xiusl/glog/glog"
)

var logger glog.GLogger

func main() {
	fmt.Println("This is my Logger")
	// logger := glog.NewLogger("debug")
	// logger := glog.NewFileLogger("debug", "./", "web.log")
	logConfig := &glog.GLoggerConfig{
		Level: "debug",
		Type:  "console",
	}
	logger = glog.NewLogger(logConfig)

	for {
		logger.Debug("log debug")
		logger.Info("log info")
		logger.Warning("log warning")

		logger.Debug("log debug")
		user := "Tom"
		logger.Info("log info user: %v", user)
		logger.Warning("log warning")
		logger.Error("this is error")
		logger.Fatal("this is Fatal")
		logger.Error("this is error")
		logger.Fatal("this is Fatal")
		logger.Error("this is error")
		logger.Fatal("this is Fatal")
		logger.Error("this is error")
		logger.Fatal("this is Fatal")
		logger.Error("this is error")
		logger.Fatal("this is Fatal")

		// time.Sleep(1 * time.Second)
	}
}
