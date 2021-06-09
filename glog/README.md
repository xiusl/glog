# glog base

# Feature

- 支持向不同位置输出日志（os.Stdout，file，db等）
- 日志级别控制
  - Debug/Trace/Info/Warning/Error/Fatal
  - 可控制不同级别的输出
- 日志格式
  - 时间、行号、文件名、函数名、日志级别、日志信息
- 日志切割（文件）
- 异步写日志（文件）

# Sync

启动多个 goroutine 写入日志，其实效率可能并没有一个 goroutine 高。
判断是因为使用多个 goroutine 时，需要在日志写入时加锁，效率可能就没有那么好了。
```go
for i := 0; i < 5; i++ {
		go f.WriteLogWorker()
}
```
