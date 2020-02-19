package main

import (
	"github.com/astaxie/beego"
	_ "kuaidian-app/routers"
)

func main() {
	beego.Run("0.0.0.0:8080")
}

