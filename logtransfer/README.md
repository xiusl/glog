# logtransfer

运输 kafka 中的消息到 ElasticSearch 中

需要单独运行

```
nohup go run logtransfer/cmd/main.go >> transfer.log 2>&1 &
```