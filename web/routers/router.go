package routers

import (
	"github.com/astaxie/beego"
	"github.com/xiusl/glog/web/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/add", &controllers.AddController{}, "get:GetFunc;post:PostFunc")
	beego.Router("/edit", &controllers.EditController{}, "get:GetFunc;post:PostFunc")
}
