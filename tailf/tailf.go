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
	t     *tail.Tail
	path  string
	topic string
}

func NewTailManager(configs []logagent.LogConfig) (*TailManager, error) {
	tm := &TailManager{
		msgChan: make(chan *TextMessage, 100),
	}
	var arr []*TailServer
	for _, c := range configs {
		t, err := newTailServer(c.Path, c.Topic)
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

func newTailServer(filename, topic string) (*TailServer, error) {
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
		t:     t,
		path:  filename,
		topic: topic,
	}, nil
}

func (m *TailManager) StartTailf(ts *TailServer) {
	log.Printf("Tail Server Start path:%v.\n", ts.path)
	for {
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
