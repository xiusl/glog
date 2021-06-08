# gLog 
 
logging service

```
               🖥🖥  <--------
                |            |
               conf          |
                |         Kibaba <---
                v                   |
               etcd                 |   
                |              ElasticSearch
                |                   ^
    |--watch-----------watch---|     |             
    v                          v     |
LogAgent  ---> Kafaka ---> tranfer --|
                                     |
                           Storm  <--|
                                     |
                          Hodoop  <--|

```

LogAgent: 日志收集客户端，用于收集服务器上的日志。
Kafaka: 高吞吐量的分布式队列。
ElasticSearch: 开源搜索引擎，提供基于 HTTP RESTFull 的 web 接口。简称ES
Kibaba: 开源 ES 数据分析和可视化工具。
Hodoop: 分布式计算框架，可以对海量数据进行分布式处理。
Storm: 开源分布式实时计算系统。