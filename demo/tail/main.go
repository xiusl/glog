package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hpcloud/tail"
)

func main() {
	fileName := "./a.log"
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		log.Fatalf("tail TailFile error: %v.\n", err)
	}
	var (
		msg *tail.Line
		ok  bool
	)

	for {
		msg, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("msg:", msg.Text)
	}
}
