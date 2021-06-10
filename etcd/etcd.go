package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/xiusl/glog/logagent"
	"go.etcd.io/etcd/clientv3"
)

var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	client, err = clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * time.Second,
		},
	)
	if err != nil {
		log.Printf("Etcd client new fail, err: %v.\n", err)
		return err
	}

	log.Println("Etcd connect success.")
	return
}

func GetConfigInfo(key string) (configs []logagent.LogConfig, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()
	resp, err := client.Get(ctx, key)
	if err != nil {
		log.Printf("Etcd get %v fail, err: %v.\n", key, err)
		return nil, err
	}
	for _, v := range resp.Kvs {
		if string(v.Key) == key {
			err = json.Unmarshal(v.Value, &configs)
			if err != nil {
				log.Printf("Config unmarshal faile error: %v.\n", err)
				continue
			}
			log.Printf("Etcd get `%v` succes\n", key)
			return
		}
	}
	err = errors.New(fmt.Sprintf("Etcd get %v fail, err: not found.\n", key))
	return
}
