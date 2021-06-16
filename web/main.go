package main

import (
	"github.com/astaxie/beego"
	"github.com/xiusl/glog/web/etcd"
	_ "github.com/xiusl/glog/web/routers"
)

const etcdAddr = "http://192.144.171.238:2379"

func main() {

	etcd.Init([]string{etcdAddr})

	beego.Run(":8083")
}
