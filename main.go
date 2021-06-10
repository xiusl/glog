package main

import (
	"log"
	"time"

	"github.com/xiusl/glog/kafka"
	"github.com/xiusl/glog/tailf"
)

func main() {
	err := kafka.Init([]string{"127.0.0.1:9092"})
	if err != nil {
		log.Fatalf("kafka init error: %v.\n", err)
	}

	err = tailf.Init("./a.log")
	if err != nil {
		log.Fatalf("tailf init error: %v.\n", err)
	}

	for {
		select {
		case line := <-tailf.ReadLine():
			kafka.SendMessageToKafka("web", line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}
