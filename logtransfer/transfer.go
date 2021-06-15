package logtransfer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/xiusl/glog/logging"
)

type TransferServer struct {
	consumer sarama.Consumer
	msgChan  chan *TransMessage
}

type TransMessage struct {
	Message string
	Topic   string
}

func NewTransferServer(address []string) (*TransferServer, error) {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		logging.Error("Consumer create fail %v.", err)
		return nil, err
	}
	return &TransferServer{
		consumer: consumer,
	}, nil
}

var (
	consumer      sarama.Consumer
	MsgChan       chan *TransMessage
	mu            sync.Mutex
	topics        map[string]int
	wg            sync.WaitGroup
	consumerGroup sarama.ConsumerGroup
	myc           *myConsumer
)

func Init(addrs []string) (err error) {
	consumer, err = sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		logging.Error("Consumer create fail %v.", err)
		return
	}
	conf := sarama.NewConfig()
	consumerGroup, err = sarama.NewConsumerGroup(addrs, "abc", conf)
	if err != nil {
		logging.Error("Consumer create fail %v.", err)
		return
	}

	mu = sync.Mutex{}
	wg = sync.WaitGroup{}
	MsgChan = make(chan *TransMessage, 100)
	topics = make(map[string]int)
	myc = &myConsumer{}
	fmt.Println("Kafka consumer init success.")
	return
}

func ConsumeMessageFromKafka(topic string) error {
	partitionList, err := consumer.Partitions("a_log")
	if err != nil {
		log.Panicf("sarama Partitions get fail err %v.\n", err)
	}
	log.Printf("partition count: %d", len(partitionList))
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("a_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("sarama ConsumePartition get fail err %v.\n", err)
			continue
		}
		defer pc.AsyncClose()

		go func() {
			for {
				fmt.Println("111111")
				for msg := range pc.Messages() {
					fmt.Printf("partition: %d, offset: %d, key: %v, value:%v.\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
					transMsg := &TransMessage{
						Message: string(msg.Value),
						Topic:   topic,
					}
					mu.Lock()
					defer mu.Unlock()
					MsgChan <- transMsg
					// select {
					// case
					// 	fmt.Println("xieruxiaoxi")
					// }

				}
			}
		}()

	}
	fmt.Println("return")
	select {}
	// return nil
	//
}
func ConsumeMessage(topics []string) {
	// consumer.ConsumePartition()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := consumerGroup.Consume(ctx, topics, myc)
	if err != nil {
		log.Printf("ConsumeMessage Error: %v.\n", err)
	}
}

type myConsumer struct {
}

func (c *myConsumer) Setup(s sarama.ConsumerGroupSession) error {
	log.Print("myConsumer Setup")
	return nil
}
func (c *myConsumer) Cleanup(s sarama.ConsumerGroupSession) error {
	log.Print("myConsumer Cleanup")
	return nil
}
func (c *myConsumer) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		key := string(message.Key)
		val := string(message.Value)
		log.Printf("%s-%s", key, val)

		transMsg := &TransMessage{
			Message: string(message.Value),
			Topic:   claim.Topic(),
		}

		select {
		case MsgChan <- transMsg:
			log.Print("123")
		}

	}

	return nil
}

func ConsumeMessageFromKafka1(topic string) error {
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Panicf("sarama Partitions get fail err %v.\n", err)
		return err
	}
	if _, exist := topics[topic]; exist {
		log.Printf("sarama Partitions topic `%v` exist.\n", topic)
		return errors.New("Topic exist")
	}
	topics[topic] = 1
	// wg.Add(1)

	topic = "a_log"

	log.Printf("ConsumeMessageFromKafka topic: `%v`, ps: %d Starting... ", topic, len(partitionList))
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("sarama ConsumePartition get fail err %v.\n", err)
			continue
		}
		defer pc.AsyncClose()
		// wg.Add(1)
		// go
		go func() {
			log.Printf("read sssssss： ")
			log.Print(len(pc.Messages()))
			for {
				for msg := range pc.Messages() {
					fmt.Printf("partition: %d, offset: %d, key: %v, value:%v.\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				}
				// 	select {
				// 	case msg := <-pc.Messages():
				// 		if msg == nil {
				// 			log.Printf("msg nil")
				// 			time.Sleep(500 * time.Millisecond)
				// 			continue
				// 		}
				// 		log.Printf("partition: %d, offset: %d, key: %v, value:%v.\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				// 	}
			}
			// log.Printf("read eeeeeee")
			// wg.Done()
		}()
		// wg.Wait()
	}
	return nil
}

/*
// for msg := range pc.Messages() {
				// 	log.Printf("partition: %d, offset: %d, key: %v, value:%v.\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				// 	transMsg := &TransMessage{
				// 		Message: string(msg.Value),
				// 		Topic:   topic,
				// 	}
				// 	log.Print(MsgChan)
				// 	log.Print(transMsg)
				// 	// mu.Lock()
				// 	// defer mu.Unlock()
				// 	select {
				// 	case MsgChan <- transMsg:
				// 		fmt.Println("xieruxiaoxi")
				// 	}
				// 	log.Printf("read okokokok")
				// }
*/
