package es

import (
	"context"
	"log"

	"github.com/olivere/elastic/v7"
)

var (
	esClient *elastic.Client
)

func Init(addr string) (err error) {
	esClient, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	if err != nil {
		log.Printf("ES Client init Fail, error: %v.\n", err)
		return
	}
	log.Printf("ES Client init success.\n")
	return
}

func SendToEs(topic, message string) {
	msg := struct {
		Topic   string
		Message string
	}{
		Topic:   topic,
		Message: message,
	}

	_, err := esClient.Index().Index(topic).BodyJson(&msg).Do(context.Background())

	if err != nil {
		log.Printf("SendToEs Failure, error: %v.\n", err)
		return
	}
	log.Println("SendToEs Success.")
	return
}
