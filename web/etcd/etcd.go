package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var NotFoundErr = errors.New("Not Found")

type EtcdKeyV struct {
	Key     string
	Value   string
	Configs []EtcdLogConfig
}

type EtcdLogConfig struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
	Key   string `json:"key"`
}

func (kv *EtcdKeyV) Contains(conf *EtcdLogConfig) bool {
	for _, c := range kv.Configs {
		if c.Path == conf.Path {
			log.Printf("c.Path: %v, conf.Path: %v.\n", c.Path, conf.Path)
			return true
		}
	}
	return false
}

func (kv *EtcdKeyV) ConfByPath(path string) *EtcdLogConfig {
	for _, c := range kv.Configs {
		if c.Path == path {
			return &c
		}
	}
	return nil
}

var (
	etcdServer *EtcdServer
)

func Init(endpoints []string) (err error) {
	etcdServer, err = NewEtcdServer(endpoints)
	return
}

type EtcdServer struct {
	client *clientv3.Client
}

func NewEtcdServer(endpoints []string) (*EtcdServer, error) {
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
		client: client,
	}, err
}

func GetAllKeys() ([]EtcdKeyV, error) {
	keyPref := ""
	resp, err := etcdServer.client.Get(context.Background(), keyPref, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var kvs []EtcdKeyV
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		value := string(kv.Value)

		var configs []EtcdLogConfig
		err = json.Unmarshal(kv.Value, &configs)
		if err != nil {
			log.Printf("Unmarshal Config failure key:%v, error: %v.\n", key, err)
			continue
		}

		tmp := EtcdKeyV{
			Key:     key,
			Value:   value,
			Configs: configs,
		}
		kvs = append(kvs, tmp)
	}
	return kvs, nil
}

func GetByKey(key string) (keyv *EtcdKeyV, err error) {
	resp, err := etcdServer.client.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		value := string(kv.Value)

		var configs []EtcdLogConfig
		err = json.Unmarshal(kv.Value, &configs)
		if err != nil {
			log.Printf("Unmarshal Config failure key:%v, error: %v.\n", key, err)
			continue
		}

		keyv = &EtcdKeyV{
			Key:     key,
			Value:   value,
			Configs: configs,
		}
		return
	}
	return nil, NotFoundErr
}

func Save(kv *EtcdKeyV) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()

	data, err := json.Marshal(kv.Configs)
	if err != nil {
		return err
	}

	_, err = etcdServer.client.Put(ctx, kv.Key, string(data))
	if err != nil {
		return err
	}
	return nil
}

func Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()

	_, err := etcdServer.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
