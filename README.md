# gLog 

Massive log collection system.  Practice project of Golang.

海量日志收集系统，Go 语言练习项目。



## 主要架构

<img src="http://pp.video.sleen.top/uPic/blog/UML%20%E5%9B%BE%20(2)-FAGOAL.jpg!blog360" alt="架构图"  />

### 名词

- LogAgent: 日志收集客户端，用于收集服务器上的日志。

- Kafka: 高吞吐量的分布式队列。

- ElasticSearch: 开源搜索引擎，提供基于 HTTP RESTFull 的 web 接口。简称ES

- Kibaba: 开源 ES 数据分析和可视化工具。

- Hodoop: 分布式计算框架，可以对海量数据进行分布式处理。

- Storm: 开源分布式实时计算系统。

### 介绍

1. 每台应用服务器上分别部署 LogAgent，每个 LogAgent 从 Etcd 中根据自己的 IP 获取指定的配置。

   ```
   // Etcd 中存储的配置
   "/logagent/192.168.0.2/config": {"path":"/logs/web.log", "topic": "web_log"}
   "/logagent/192.168.0.3/config": {"path":"/logs/im.log", "topic": "im_log"}
   ```

   <img src="http://pp.video.sleen.top/uPic/blog/%E6%B5%81%E7%A8%8B%E5%9B%BE-RiJobR.png!blog360" alt="流程图" />

2. LogAgent 实时监听 Etcd 中对应 key 的变化，实现热更新。

3. LogAgent 根据配置信息，读取对应的日志文件发送到指定 topic 的 kafka 中。

4. 收集在 Kafka 中的日志再通过 LogTransfer 传送到 ES、Hadoop 等。

5. 在 ElasticSearch 中再结合 Kibana 进行可视化的数据分析。

6. 对应的也可以使用 Hadoop 和 storm 进行大数据分析。

## 项目结构

```
- glog/
  - demo/        各种三方库的简单应用
  - es/          发送消息到 ES
  - etcd/        读取和监听配置信息
  - kafka/       kafka 生产者封装
  - logging/     日志库
  - transfer/    从 Kafka 读取消息
    - cmd/            发送消息到指定的位置（ES，Headoop等）       
    - transfer.go     对 Kafka 读取消息的封装  
  - tailf/       监听读取日志文件内容
  - main.go
```

## 应用

对 nginx 日志的收集

以哩嗑 https://ins.sleen.top 项目为例

1. 部署 etcd (可以参考 References/etcd.md 或自行 Google)，我这里是直接部署在我的服务器上 192.144.171.238:2379
2. 部署 kafka (可以参考 References/kafka.md 或自行 Google)，同样我部署在了 192.144.171.238:9092
3. 部署 glog/web/main.go，这个主要是给 etcd 提供了一个可视化的 key-value 编辑工具，便于修改配置信息，这个也放在了 192.144.171.238:8083，完成后添加以下 key-value
  ```
  /bd/logagent/config/122.112.235.92 : [{"topic":"nginx_log", "path":"/var/log/nginx/access.log"}]
  ```
4. 修改 config 下的配置信息，并在 ins.sleen.top 的服务器上运行 glog/main.go，这样 access.log 的内容就会发送到 kafka 中

5. 部署 ElasticSearch，(可以参考 demo/elasticsearch/readme.md 或自行 Google)，依旧是放在 192.144.171.238:9200

6. 在一台服务器上运行 glog/logtransfer/cmd/main.go ，读取 kafka 消息传输到 Es 中

7. 在本地或其他服务器上安装 kibana，配置 Es host，设置索引，即可对日志信息进行分析 

   

基本效果，当 ins.sleen.top 被访问时，日志信息会收集起来

![20210616141012](http://pp.video.sleen.top/uPic/blog/20210616141012-JhsMKq.jpg)




## 参考

- [GO语言全栈工程师 287P完结](https://www.bilibili.com/video/BV1FV411r7m8)
- [Golang实战之海量日志收集系统](https://blog.csdn.net/qq_43442524/article/details/105023724)

