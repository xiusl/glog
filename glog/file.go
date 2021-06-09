package glog

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type FileLogger struct {
	Level    Level
	FilePath string
	FileName string
	MaxSize  int
	file     *os.File
	errFile  *os.File
}

func NewFileLogger(levelStr string, filePath, fileName string) *FileLogger {
	level, err := parseLevel(levelStr)
	if err != nil {
		panic(err)
	}

	f := &FileLogger{
		Level:    level,
		FilePath: filePath,
		FileName: fileName,
	}

	err = f.initFile()
	if err != nil {
		panic(err)
	}

	return f
}

func (f *FileLogger) initFile() error {
	fullPath := path.Join(f.FilePath, f.FileName)
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("FileLogger: Failed to open log file, err:%v.\n", err)
		return err
	}
	f.file = file

	errFileName := strings.Replace(f.FileName, ".log", "_err.log", 1)
	fullPath = path.Join(f.FilePath, errFileName)
	errFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("FileLogger: Failed to open log file, err:%v.\n", err)
		return err
	}
	f.errFile = errFile
	return nil
}

func (f *FileLogger) enable(level Level) bool {
	return f.Level <= level
}

func (f *FileLogger) Log(level Level, format string, arg ...interface{}) {
	if !f.enable(level) {
		return
	}
	msg := fmt.Sprintf(format, arg...)
	timeStr := time.Now().Format(timeFromatStr)
	filePath, fileName, lineNo := getInfo(3)

	logStr := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", timeStr, level.ToString(), fileName, filePath, lineNo, msg)
	f.file.Write([]byte(logStr))

	if level >= ERROR {
		_, err := f.errFile.Write([]byte(logStr))
		if err != nil {
			fmt.Printf("Errfile: Failed to Wirte, err: %v.\n", err)
		}
	}
}

func (f *FileLogger) Debug(format string, arg ...interface{}) {
	f.Log(DEBUG, format, arg...)
}

func (f *FileLogger) Info(format string, arg ...interface{}) {
	f.Log(INFO, format, arg...)
}

func (f *FileLogger) Warning(format string, arg ...interface{}) {
	f.Log(WARNING, format, arg...)
}

func (f *FileLogger) Error(format string, arg ...interface{}) {
	f.Log(ERROR, format, arg...)
}

func (f *FileLogger) Fatal(format string, arg ...interface{}) {
	f.Log(FATAL, format, arg...)
}
