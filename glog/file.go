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
	MaxSize  int64
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
		MaxSize:  10 * 1024 * 1024,
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

	if f.checkFileSize(f.file) {
		if file, ok := f.splitFile(f.file); ok {
			f.file = file
		}
	}

	logStr := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", timeStr, level.ToString(), fileName, filePath, lineNo, msg)
	_, err := f.file.Write([]byte(logStr))
	if err != nil {
		fmt.Printf("FileLog: Failed to Wirte log, err: %v.\n", err)
	}
	if level >= ERROR {
		if f.checkFileSize(f.errFile) {
			if file, ok := f.splitFile(f.errFile); ok {
				f.errFile = file
			}
		}

		_, err := f.errFile.Write([]byte(logStr))
		if err != nil {
			fmt.Printf("FileLog: Failed to Wirte log, err: %v.\n", err)
		}
	}
}

func (f *FileLogger) checkFileSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Failed to check file size, err: %v\n", err)
		return false
	}
	return (fileInfo.Size() > f.MaxSize)
}

func (f *FileLogger) splitFile(file *os.File) (*os.File, bool) {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("SplitFile: Failed to get file info, err: %v\n", err)
		return nil, false
	}
	fileName := fileInfo.Name()
	bkFormat := fmt.Sprintf("%s.log", time.Now().Format("_20060102150405000"))
	backupName := strings.Replace(fileName, ".log", bkFormat, 1)
	backupName = path.Join(f.FilePath, backupName)
	logName := path.Join(f.FilePath, fileName)

	// close old file
	file.Close()
	// backup
	err = os.Rename(logName, backupName)
	if err != nil {
		fmt.Printf("SplitFile: Failed to backup the old log file, err: %v", err)
		return nil, false
	}

	newFile, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("SplitFile: Failed to create new log file, err: %v", err)
		return nil, false
	}
	return newFile, true
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
