package routers

import (
	"github.com/astaxie/beego"
	"kuaidian-app/controllers"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
