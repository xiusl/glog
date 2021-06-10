# tail

## 安装
```
go get github.com/hpcloud/tail
```

## 简介

实时监控文件新增内容，类似于 `linux` 的 `tail -f` 命令。

## 使用
`main.go`

``` go
// 光标位置
// Offset: 偏移量
// Whence: 0 => 从头，1 => 当前，2 => 结尾 
// os.Seek
tail.SeekInfo{Offset: 0, Whence: 2}

// 轮询文件修改，而不是使用 inotify
poll 
```