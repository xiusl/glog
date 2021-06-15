package main

import (
	"github.com/astaxie/beego"
	"github.com/xiusl/glog/web/etcd"
	_ "github.com/xiusl/glog/web/routers"
)

func main() {

	etcd.Init([]string{"http://127.0.0.1:2379"})

	beego.Run(":8083")
}
