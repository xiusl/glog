package main

import (
	"fmt"
	"time"

	"github.com/xiusl/glog/glog"
)

func main() {
	fmt.Println("This is my Logger")
	// logger := glog.NewLogger("debug")
	logger := glog.NewFileLogger("debug", "./", "web.log")

	n := time.Now()

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
		if n.Add(time.Minute * 1).Before(time.Now()) {
			break
		}
	}
}
