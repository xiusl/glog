package main

import (
	"log"

	"github.com/xiusl/glog/etcd"
	"github.com/xiusl/glog/kafka"
	"github.com/xiusl/glog/setting"
	"github.com/xiusl/glog/tailf"
)

func main() {
	err := kafka.Init([]string{"127.0.0.1:9092"})
	if err != nil {
		log.Fatalf("kafka init error: %v.\n", err)
	}

	err = etcd.Init([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatalf("etcd init error: %v.\n", err)
	}

	// 从 etcd 获取 config
	configs, err := etcd.GetConfigInfo(setting.EtcdKey)
	if err != nil {
		log.Fatalf("etcd GetConfigInfo error: %v.\n", err)
	}

	// 初始化 tailMgr
	tm, err := tailf.NewTailManager(configs)
	if err != nil {
		log.Fatalf("tail servers start error: %v.\n", err)
	}

	for {
		msg := tm.ReadMessage()
		if msg != nil {
			err = kafka.SendMessageToKafka(msg.Topic, msg.Message)
			if err != nil {
				log.Printf("kafka SendMessage error: %v.\n", err)
			}
		}
	}
}
