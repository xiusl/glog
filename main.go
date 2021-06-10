package main

import (
	"log"

	"github.com/xiusl/glog/etcd"
	"github.com/xiusl/glog/kafka"
	"github.com/xiusl/glog/logagent"
	"github.com/xiusl/glog/setting"
	"github.com/xiusl/glog/tailf"
)

func main() {
	err := kafka.Init([]string{"127.0.0.1:9092"})
	if err != nil {
		log.Fatalf("kafka init error: %v.\n", err)
	}

	// 可以根据每台机器的IP来生成不同的key
	keys := []string{
		setting.EtcdKey,
		"/bd/logagent/config/0.0.0.2",
		"/bd/logagent/config/0.0.0.3",
		"/bd/logagent/config/0.0.0.4",
	}

	etcdServer, err := etcd.NewEtcdServer([]string{"127.0.0.1:2379"}, keys)
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
		// for {
		select {
		case wch := <-etcdSv.WatchChan:
			if wch.Type == "del" {
				tailMgr.StopTail(wch.Key)
			} else if wch.Type == "put" {
				tailMgr.UpdateConfig(wch.Key, wch.Confs)
			}
		default:
		}
		// }
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
