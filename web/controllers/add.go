package controllers

import (
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/xiusl/glog/web/etcd"
)

type AddController struct {
	beego.Controller
}

func (c *AddController) GetFunc() {
	flash := beego.ReadFromRequest(&c.Controller)
	if n, ok := flash.Data["notice"]; ok {
		log.Print(n)
	}
	c.TplName = "add.tpl"
}

func (c *AddController) PostFunc() {
	flash := beego.NewFlash()

	key := c.GetString("key")
	topic := c.GetString("topic")
	path := c.GetString("path")

	kv, _ := etcd.GetByKey(key)
	if kv == nil {
		kv = &etcd.EtcdKeyV{}
		kv.Key = key
	}

	conf := etcd.EtcdLogConfig{
		Topic: topic,
		Path:  path,
		Key:   key,
	}
	if kv.Contains(&conf) {
		log.Printf("Conf did existed! key: %v.\n", key)
		flash.Notice("Conf did existed!")
		flash.Store(&c.Controller)
		c.Redirect("/add", 302)
		return
	}
	kv.Configs = append(kv.Configs, conf)

	err := etcd.Save(kv)
	if err != nil {
		log.Printf("Conf save failure! error: %v.\n", err)
		flash.Notice("Conf save failure!")
		flash.Store(&c.Controller)
		c.Redirect("/add", 302)
		return
	}

	c.Redirect("/", 302)
}

type EditController struct {
	beego.Controller
}

func (c *EditController) PostFunc() {
	flash := beego.NewFlash()

	key := c.GetString("key")
	topic := c.GetString("topic")
	path := c.GetString("path")

	url := fmt.Sprintf("/edit?key=%v&path=%v&topic=%v", key, path, topic)

	kv, _ := etcd.GetByKey(key)
	if kv == nil {
		flash.Notice("Conf did existed!")
		flash.Store(&c.Controller)
		c.Redirect(url, 302)
		return
	}

	var tmp []etcd.EtcdLogConfig
	for _, c := range kv.Configs {
		if c.Path == path && c.Topic == topic {
			continue
		}
		tmp = append(tmp, c)
	}
	kv.Configs = tmp

	if len(tmp) == 0 {
		err := etcd.Delete(key)
		if err != nil {
			log.Printf("Conf save failure! error: %v.\n", err)
			flash.Notice("Conf save failure!")
			flash.Store(&c.Controller)
			c.Redirect(url, 302)
			return
		}
	} else {
		err := etcd.Save(kv)
		if err != nil {
			log.Printf("Conf save failure! error: %v.\n", err)
			flash.Notice("Conf save failure!")
			flash.Store(&c.Controller)
			c.Redirect(url, 302)
			return
		}
	}

	c.Redirect("/", 302)
}

func (c *EditController) GetFunc() {
	flash := beego.ReadFromRequest(&c.Controller)
	if n, ok := flash.Data["notice"]; ok {
		log.Print(n)
	}
	c.TplName = "edit.tpl"
}
