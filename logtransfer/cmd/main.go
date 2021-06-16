package main

import (
	"log"

	"github.com/xiusl/glog/config"
	"github.com/xiusl/glog/es"
	"github.com/xiusl/glog/etcd"
	"github.com/xiusl/glog/logagent"
	"github.com/xiusl/glog/logging"
	"github.com/xiusl/glog/logtransfer"
)

type Array []string

func main() {
	logging.Info("LogTransfer Starting...")

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config load error: %v.\n", err)
	}

	// 从 Etcd 获取 transfer 配置
	keys := []string{conf.EtcdKey}
	etcdAddrs := []string{conf.EtcdAddr}
	etcdServer, err := etcd.NewEtcdServer(etcdAddrs, keys)
	if err != nil {
		panic(err)
	}

	var configs []logagent.LogConfig
	for _, key := range keys {
		tmpConfigs, err := etcdServer.GetConfigInfo(key)
		if err != nil {
			log.Printf("etcd GetConfigInfo error: %v.\n", err)
			continue
		}
		configs = append(configs, tmpConfigs...)
	}

	kafkaAddrs := []string{conf.KafkaAddr}
	// trans, err := logtransfer.NewTransferServer(kafkaAddrs)
	// if err != nil {
	// 	panic(err)
	// }
	logtransfer.Init(kafkaAddrs)

	// 初始化Es
	es.Init(conf.EsAddr)

	go func() {
		for {
			select {
			case msg := <-logtransfer.MsgChan:
				// 转发到 Es
				log.Printf("转发至Es, %v", msg.Message)
				es.SendToEs(msg.Topic, msg.Message)
			}
		}
	}()

	var topics Array
	for _, conf := range configs {
		if topics.contains(conf.Topic) {
			continue
		}
		topics = append(topics, conf.Topic)
		// logtransfer.ConsumeMessageFromKafka(conf.Topic)
		// break
	}
	log.Printf("Topics %v.\n", topics)

	logtransfer.ConsumeMessage(topics)

	// select {}
	// 消费 Kafka 消息

}

func (arr Array) contains(item string) bool {
	for _, str := range arr {
		if item == str {
			return true
		}
	}
	return false
}
