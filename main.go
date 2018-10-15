package main

import (
	"github.com/astaxie/beego"
	_ "questionaire/routers"
)

func main() {
	beego.BConfig.WebConfig.EnableDocs = true
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
