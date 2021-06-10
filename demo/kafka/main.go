package main

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		log.Panicf("sarama NewConsumer init fail err %v\n.", err)
	}
	partitionList, err := consumer.Partitions("web")
	if err != nil {
		log.Panicf("sarama Partitions get fail err %v\n.", err)
	}
	log.Printf("partition count: %d", len(partitionList))
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("web", int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("sarama ConsumePartition get fail err %v\n.", err)
			continue
		}
		defer pc.AsyncClose()
		for {
			go func() {
				for msg := range pc.Messages() {
					fmt.Printf("partition: %d, offset: %d, key: %v, value:%v.\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				}
			}()
		}
	}
}

func sendToKafka() {
	config := sarama.NewConfig()
	// 发送完数据需要 leader 和 follow 都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个 partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的消息将在 success chan 返回
	config.Producer.Return.Successes = true

	msg := &sarama.ProducerMessage{}
	msg.Topic = "web"
	msg.Value = sarama.StringEncoder("this is a message")

	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Panicf("sarama NewSyncProducer init fail err %v\n.", err)
	}
	defer client.Close()

	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		log.Panicf("sarama SendMessages fail err %v\n.", err)
	}
	fmt.Printf("pid: %v, offset: %v", pid, offset)
}
