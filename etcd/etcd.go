package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/xiusl/glog/logagent"
	"go.etcd.io/etcd/clientv3"
)

type EtcdServer struct {
	client    *clientv3.Client
	keys      []string
	WatchChan chan *EtcdWatchMessage
}

type EtcdWatchMessage struct {
	// del, put
	Type  string
	Key   string
	Confs []logagent.LogConfig
}

func NewEtcdServer(endpoints, keys []string) (*EtcdServer, error) {
	client, err := clientv3.New(
		clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		},
	)
	if err != nil {
		log.Printf("Etcd client new fail, err: %v.\n", err)
		return nil, err
	}

	log.Println("Etcd connect success.")
	return &EtcdServer{
		client:    client,
		keys:      keys,
		WatchChan: make(chan *EtcdWatchMessage, 10),
	}, err
}

func (s *EtcdServer) GetConfigInfo(key string) (configs []logagent.LogConfig, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()
	resp, err := s.client.Get(ctx, key)
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
			log.Printf("Etcd get `%v` success.\n", key)
			return
		}
	}
	err = errors.New(fmt.Sprintf("Etcd get %v fail, err: not found.\n", key))
	return
}

func (s *EtcdServer) WatchKeys() {
	for _, key := range s.keys {
		go s.watchKey(key)
	}
}

func (s *EtcdServer) watchKey(key string) {
	wch := s.client.Watch(context.Background(), key)
	for resp := range wch {
		for _, ev := range resp.Events {
			// key 被删除了
			if ev.Type == mvccpb.DELETE {
				// 通知 tail Mgr 删除对当前 key 下文件的读取
				wt := &EtcdWatchMessage{
					Type: "del",
					Key:  key,
				}
				log.Printf("EtcdServer watchKey %v delete", key)
				select {
				case s.WatchChan <- wt:
				}

			}
			if ev.Type == mvccpb.PUT {
				var confs []logagent.LogConfig
				err := json.Unmarshal(ev.Kv.Value, &confs)
				if err != nil {
					log.Printf("Etcd WatchKey %v, value Unmarshal faile, err:%v.\n", key, err)
					continue
				}
				log.Printf("EtcdServer watchKey %v put", key)
				// 通知 tail Mgr 指定 key 的 confs 更新了
				wt := &EtcdWatchMessage{
					Type:  "put",
					Key:   key,
					Confs: confs,
				}
				select {
				case s.WatchChan <- wt:
				}
			}
		}
	}
}
