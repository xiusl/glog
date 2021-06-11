# gLog 

Massive log collection system.  Practice project of Golang.

海量日志收集系统，Go 语言练习项目。



## 主要架构

<img src="http://pp.video.sleen.top/uPic/blog/UML%20%E5%9B%BE%20(2)-FAGOAL.jpg" alt="架构图" style="zoom:40%;" />

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

   <img src="http://pp.video.sleen.top/uPic/blog/%E6%B5%81%E7%A8%8B%E5%9B%BE-RiJobR.png" alt="流程图" style="zoom:50%;" />

2. LogAgent 实时监听 Etcd 中对应 key 的变化，实现热更新。

3. LogAgent 根据配置信息，读取对应的日志文件发送到指定 topic 的 kafka 中。

4. 收集在 Kafka 中的日志再通过 LogTransfer 传送到 ES、Hadoop 等。

5. 在 ElasticSearch 中再结合 Kibana 进行可视化的数据分析。

6. 对应的也可以使用 Hadoop 和 storm 进行大数据分析。

## 参考

- [GO语言全栈工程师 287P完结](https://www.bilibili.com/video/BV1FV411r7m8)
- [Golang实战之海量日志收集系统](https://blog.csdn.net/qq_43442524/article/details/105023724)

