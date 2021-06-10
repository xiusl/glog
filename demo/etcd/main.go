package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
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

}

func etcdLease(cli *clientv3.Client) {
	// lease
	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		log.Fatalf("etcd client Grant new fail, err: %v.\n", err)
	}

	// 5s 后，token 将被移除
	_, err = cli.Put(context.Background(), "token", "abc", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatalf("etcd cli put fail, err: %v\n", err)
	}

	for {
		resp, err := cli.Get(context.Background(), "token")
		if err != nil {
			log.Fatalf("etcd cli get fail, err: %v\n", err)
		}
		if resp.Count == 0 {
			fmt.Println("token removed")
			break
		}
		for _, v := range resp.Kvs {
			fmt.Printf("%s:%s\n", v.Key, v.Value)
		}
		time.Sleep(1 * time.Second)
	}
}
func etcdWatch(cli *clientv3.Client) {
	// watch
	wat := cli.Watch(context.Background(), "name")
	for resp := range wat {
		for _, ev := range resp.Events {
			fmt.Printf("Type: %v, Key: %v, Value: %v.\n", ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
		}
	}
}

func etcdPutGet() {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()

	// put
	_, err = cli.Put(ctx, "name", "tom")
	if err != nil {
		log.Fatalf("etcd cli put fail, err: %v\n", err)
	}

	// get
	resp, err := cli.Get(ctx, "name")
	if err != nil {
		log.Fatalf("etcd cli get fail, err: %v\n", err)
	}

	for _, v := range resp.Kvs {
		fmt.Printf("%s:%s\n", v.Key, v.Value)
	}
}
