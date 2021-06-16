package main

import (
	"log"

	"github.com/xiusl/glog/config"
	"github.com/xiusl/glog/etcd"
	"github.com/xiusl/glog/kafka"
	"github.com/xiusl/glog/logagent"
	"github.com/xiusl/glog/tailf"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config load error: %v.\n", err)
	}

	err = kafka.Init([]string{conf.KafkaAddr})
	if err != nil {
		log.Fatalf("kafka init error: %v.\n", err)
	}

	// 可以根据每台机器的IP来生成不同的key
	keys := []string{
		conf.EtcdKey,
	}

	etcdServer, err := etcd.NewEtcdServer([]string{conf.EtcdAddr}, keys)
	if err != nil {
		log.Fatalf("etcd init error: %v.\n", err)
	}

	// 从 etcd 获取 config
	var configs []logagent.LogConfig

	for _, key := range keys {
		tmpConfigs, err := etcdServer.GetConfigInfo(key)
		if err != nil {
			log.Printf("etcd GetConfigInfo error: %v.\n", err)
			continue
		}
		configs = append(configs, tmpConfigs...)
	}

	// 监听所有的 keys
	etcdServer.WatchKeys()

	// 初始化 tailMgr
	tm, err := tailf.NewTailManager(configs)
	if err != nil {
		log.Fatalf("tail servers start error: %v.\n", err)
	}

	//
	go func(etcdSv *etcd.EtcdServer, tailMgr *tailf.TailManager) {
		select {
		case wch := <-etcdSv.WatchChan:
			if wch.Type == "del" {
				tailMgr.StopTail(wch.Key)
			} else if wch.Type == "put" {
				tailMgr.UpdateConfig(wch.Key, wch.Confs)
			}
		}
	}(etcdServer, tm)

	// 监听消息转发到 kafka
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
