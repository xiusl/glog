package tailf

import (
	"fmt"
	"log"

	"github.com/hpcloud/tail"
	"github.com/xiusl/glog/logagent"
)

type TailManager struct {
	tails   []*TailServer
	msgChan chan *TextMessage
}

type TextMessage struct {
	Message string
	Topic   string
}

type TailServer struct {
	key   string
	t     *tail.Tail
	path  string
	topic string
	stop  bool
}

func NewTailManager(configs []logagent.LogConfig) (*TailManager, error) {
	tm := &TailManager{
		msgChan: make(chan *TextMessage, 100),
	}
	var arr []*TailServer
	for _, c := range configs {
		t, err := newTailServer(c.Path, c.Topic, c.Key)
		if err != nil {
			return nil, err
		}
		arr = append(arr, t)

		//
		go tm.StartTailf(t)
	}
	tm.tails = arr
	return tm, nil
}

func newTailServer(filename, topic, key string) (*TailServer, error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Printf("tail TailFile error: %v.\n", err)
		return nil, err
	}
	return &TailServer{
		key:   key,
		t:     t,
		path:  filename,
		topic: topic,
	}, nil
}

func (m *TailManager) StartTailf(ts *TailServer) {
	log.Printf("Tail Server Start path:%v.\n", ts.path)
	for {
		if ts.stop {
			return
		}
		select {
		case line, ok := <-ts.t.Lines:
			if !ok {
				log.Printf("Tail Server read fail path:%v.\n", ts.path)

				continue
			}
			msg := &TextMessage{
				Message: line.Text,
				Topic:   ts.topic,
			}
			select {
			case m.msgChan <- msg:
			default:
			}
		default:
		}
	}
}

func (tm *TailManager) ReadMessage() *TextMessage {
	select {
	case line := <-tm.msgChan:
		return line
	default:
		return nil
	}
}

func (tm *TailManager) StopTail(key string) {
	// 一个 key 可能对应多个 tailServer
	var delIndexs []int
	for i, ts := range tm.tails {
		if ts.key == key {
			ts.stop = true
			ts.t.Stop()
			ts.t.Cleanup()
			delIndexs = append(delIndexs, i)
			log.Printf("TailManager StopTail key: %v, path: %v.\n", key, ts.path)
		}
	}
	for _, v := range delIndexs {
		tm.tails = append(tm.tails[:v], tm.tails[v+1:]...)
	}
}

func (tm *TailManager) UpdateConfig(key string, configs []logagent.LogConfig) {
	log.Printf("TailManager UpdateConfig Start key: %v.\n", key)
	var newTses []*TailServer
	for _, conf := range configs {
		exist := false
		for _, ts := range tm.tails {
			if ts.path == conf.Path {
				exist = true
			}
		}
		if exist {
			continue
		}
		t, err := newTailServer(conf.Path, conf.Topic, key)
		if err != nil {
			log.Printf("UpdateConfig newTailServer error: %v.\n", err)
			continue
		}
		log.Printf("UpdateConfig newTailServer key:%v, path:%v.\n", key, t.path)
		newTses = append(newTses, t)

		//
		go tm.StartTailf(t)
	}
	tm.tails = append(tm.tails, newTses...)
}

func (tm *TailManager) getTailServerByKey(key string) (*TailServer, int) {
	for i, v := range tm.tails {
		if v.key == key {
			return v, i
		}
	}
	return nil, -1
}
