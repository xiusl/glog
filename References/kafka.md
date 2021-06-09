# Kafka

## 简介

由 LinkedIn 开发，使用 Scala 语言，现为 Apache 开源项目。
分布式数据流平台，可以运行在单台服务器，也可以部署为集群。
提供了发布/订阅的功能。
具有高吞吐、低延迟、高容错等特点。

## 基本概念

- **Producer:** 生产者，消息的生产者，消息的入口。

- **Kafka cluster:** Kafka 集群，一台或多台服务器组成。
  - **Broker:** 部署了 Kafka 实例的服务器节点。每台服务器上有一个或多个 Kafka 实例。每个 kafka 集群内的 broker 都有一个不重复的编号。
  - **Topic:** 消息的主题，可以理解为消息的分类。数据保存在 Topic 上。每个 Broker 上可以创建多个 Topic。实际应用中通常一个业务线一个 Topic。
  - **Partition:**  Topic 的分区，每个 Topic 可以有多个分区，分区的作用是做负载，提高 kafka 的吞吐量。同一个 topic 在不同分区的数据是不重复的，partition 的表现形式就是一个一个的文件夹。
  - **Replication:** 每个分区都有好多副本，副本的作用是做备胎。当主分区 （Leader）故障时，会选择一个备胎（Follow）上位，成为 Leader。默认副本最大数量是 10 个，且副本数量不能大于 broker 的数量，follow 和 leader 绝对在不同的机器，同一机器对同一分区也只能存放一个副本（包括自己）。
  
- **Consumer:** 消费者，消息的消费方，消息的出口。

  - **Consumer Group:** 将多个消费者组成一个消费者组，在 kfaka 的设计中同一个分区的数据只能被消费者组中的某一个消费者消费。同一消费者组的消费者可以消费同一个 Topic 的不同分区的数据，这是为了提高吞吐量。

## 工作流程

![UML 图](http://pp.video.sleen.top/uPic/blog/UML%20%E5%9B%BE-Bmefxz.jpg)

1. 生产者从 kafka 集群获取分区 leader 信息
2. 生产者将消息发送给 leader
3. leader 将消息写入本地磁盘
4. follower 从 leader 拉取消息数据
5. follower 将消息写入本地磁盘后向 leader 发送 ACK
6. leader 收到所有的 follower 的 ACK 之后向生产者发送 ACK



## 选择 partition 的原则

在 kafka 中，如果某个 topic 有多个 partition，该如何选择发送到哪个 partition。

几个原则：

1. 在写入的时候可指定需要写入的 partition，如果有指定，那么写入指定的 partition。
2. 没有指定，但设置了数据的 key，则会根据 key 的值 hash 出一个 partition。
3. 没有指定，也没有设置 key，则会采用 **轮询** 的方式，即每次取一小段时间的数据写入某个 partition，下一小段写入下一个 partition。



## ACK 应答机制

producer 在向 kafka 写入消息时。可以设置参数来确定 kafka 是否收到数据，参数的值可以设置为 0，1，all。

- **0:**  代表 producer 往 kafka 发送的数据不需要等待返回，不确保消息发送成功，安全性最低，但效率高。
- **1:** 代表发送消息后，只要 leader 应答就可以发送下一条消息，只确保了 leader 发送成功。
- **all:** 代表向集群发送的数据需要所有 follower 都完成从 leader 的同步才会发送下一条，确保 leader 发送成功，且所有副本完成备份。安全性最好，但效率最低。

**需要注意**：如果向不存在的 Topic 发送数据，kafka 会自动创建 Topic，partition和replication，数量都是 1。

## Mac 安装

```shell
brew install kafka

# 后台启动
zookeeper-server-start -daemon /usr/local/etc/kafka/zookeeper.properties
kafka-server-start -daemon /usr/local/etc/kafka/server.properties

# brew 启动
brew services start kafka

# 启动
zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties
kafka-server-start /usr/local/etc/kafka/server.properties
```

## 测试

```
$ kafka-console-producer --bootstrap-server 127.0.0.1:9092 --topic "web"
> 123
> 456
```

```
$ kafka-console-consumer --bootstrap-server 127.0.0.1:9092 --topic "web"
> 123
> 456
```







