package main

import (
	"context"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

func main() {
	url := "http://127.0.0.1:9200/"
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(url))
	if err != nil {
		log.Fatal(err)
		return
	}

	msg := struct {
		Topic   string
		Message string
	}{
		Topic:   "a_log",
		Message: "knbasdkhsoaijdlskandjsjkdnihsi",
	}

	resp, err := client.Index().Index("a_log").BodyJson(&msg).Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(resp)
	return
}
