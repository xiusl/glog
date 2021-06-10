# etcd

## 简介

开源高可用的分布式 key-value 存储系统。使用 go 语言编写，主要参照了 ZooKeeper 和 Doozer。

支持 HTTP/JSON API，使用简单

支持 SSL 客户端证书认证

使用 Raft 算法保证强一致性，实现分布式


## 应用

- 配置中心
- 服务发现
- 分布式锁

## 安装

- Mac
```
brew install etcd
```

- Linux
[Install](https://etcd.io/docs/v3.4/install/)

## 集群

TODO

## 使用

```shell
# mac 中安装好后直接执行 etcd 命令
$ etcd
> ...
> ...
```

### 服务验证
使用 `etcdctl` 命令
```shell
$ etcdctl put name tom
> OK

$ etcdctl get name
> name
> tom

$ etcdctl del name
> 1
```

