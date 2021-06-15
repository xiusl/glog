package controllers

import (
	"log"

	"github.com/astaxie/beego"
	"github.com/xiusl/glog/web/etcd"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"

	data, err := etcd.GetAllKeys()
	if err != nil {
		log.Printf("Etcd get key failure, %v.\n", err)
	}
	c.Data["Kvs"] = data
	c.TplName = "index.tpl"
}
