package setting

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/xiusl/glog/logagent"
	"go.etcd.io/etcd/clientv3"
)

const (
	EtcdKey = "/bd/logagent/config/0.0.0.1"
)

func SetupConfigToEtcd() {

	key := "/bd/logagent/config/0.0.0.1"

	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * time.Second,
		},
	)
	if err != nil {
		log.Fatalf("etcd client new fail, err: %v.\n", err)
	}

	fmt.Println("connect to etcd...")
	defer cli.Close()

	var arr []logagent.LogConfig
	arr = append(arr, logagent.LogConfig{
		Key:   key,
		Path:  "./logs/a.log",
		Topic: "a_log",
	})
	arr = append(arr, logagent.LogConfig{
		Key:   key,
		Path:  "./logs/b.log",
		Topic: "a_log",
	})

	data, err := json.Marshal(arr)
	if err != nil {
		log.Panicf("config json marshal fail error: %v.\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()

	//
	_, err = cli.Put(ctx, key, string(data))
	if err != nil {
		log.Fatalf("etcd cli put fail, err: %v\n", err)
	}

	// get
	resp, err := cli.Get(ctx, key)
	if err != nil {
		log.Fatalf("etcd cli get fail, err: %v\n", err)
	}

	for _, v := range resp.Kvs {
		fmt.Printf("%s:%s\n", v.Key, v.Value)
	}
}
