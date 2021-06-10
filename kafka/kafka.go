package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var (
	client sarama.SyncProducer
)

func Init(address []string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		fmt.Printf("sarama NewSyncProducer init fail err %v\n.", err)
		return
	}
	return
}

func SendMessageToKafka(topic, message string) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(message)

	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Printf("sarama SendMessages fail err %v.\n", err)
		return err
	}
	fmt.Printf("pid: %v, offset: %v.\n", pid, offset)
	return nil
}
