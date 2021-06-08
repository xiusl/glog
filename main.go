package main

import (
	"fmt"
	"time"

	"github.com/xiusl/glog/glog"
)

func main() {
	fmt.Println("This is my Logger")
	logger := glog.NewLogger("debug")

	for {
		logger.Debug("log debug")
		logger.Info("log info")
		logger.Warning("log warning")
		time.Sleep(2 * time.Second)

		logger.Debug("log debug")
		user := "Tom"
		logger.Info("log info user: %v", user)
		logger.Warning("log warning")

		time.Sleep(2 * time.Second)
	}
}
