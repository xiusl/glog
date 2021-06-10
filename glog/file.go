package glog

import (
	"errors"
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
	logChan  chan *logMsg
}

type logMsg struct {
	level     Level
	msg       string
	fileName  string
	funcName  string
	timestamp string
	line      int
}

func NewFileLogger(levelStr string, filePath, fileName string) *FileLogger {
	level, err := parseLevel(levelStr)
	if err != nil {
		panic(err)
	}

	if len(filePath) == 0 {
		panic(errors.New("FileLogger need a filepath"))
	}

	if len(fileName) == 0 {
		panic(errors.New("FileLogger need a filename"))
	}

	f := &FileLogger{
		Level:    level,
		FilePath: filePath,
		FileName: fileName,
		MaxSize:  10 * 1024 * 1024,
		logChan:  make(chan *logMsg, 1000),
	}

	err = f.initFile()
	if err != nil {
		panic(err)
	}

	go f.WriteLogWorker()

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
	funcName, fileName, lineNo := getInfo(3)

	log := &logMsg{
		level:     level,
		msg:       msg,
		fileName:  fileName,
		funcName:  funcName,
		timestamp: timeStr,
		line:      lineNo,
	}

	select {
	case f.logChan <- log:
	default:
	}
}

func (f *FileLogger) WriteLogWorker() {
	for {
		select {
		case log := <-f.logChan:
			if f.checkFileSize(f.file) {
				if newFile, ok := f.splitFile(f.file); ok {
					f.file = newFile
				}
			}
			f.WriteLog(f.file, log)

			if log.level >= ERROR {
				if f.checkFileSize(f.errFile) {
					if newFile, ok := f.splitFile(f.errFile); ok {
						f.errFile = newFile
					}
				}
				f.WriteLog(f.errFile, log)
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (f *FileLogger) WriteLog(file *os.File, log *logMsg) {
	line := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", log.timestamp, log.level.ToString(),
		log.fileName, log.funcName, log.line, log.msg)

	_, err := file.Write([]byte(line))
	if err != nil {
		fmt.Printf("FileLog: Failed to Wirte log, err: %v.\n", err)
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

func (f *FileLogger) Trace(format string, arg ...interface{}) {
	f.Log(TRACE, format, arg...)
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
