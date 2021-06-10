package tailf

import (
	"fmt"

	"github.com/hpcloud/tail"
)

var tailf *tail.Tail

func Init(filename string) (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tailf, err = tail.TailFile(filename, config)
	if err != nil {
		fmt.Printf("tail TailFile error: %v.\n", err)
		return err
	}
	return nil
}

func ReadLine() <-chan *tail.Line {
	return tailf.Lines
}
