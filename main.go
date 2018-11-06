package main

import (
	_ "questionaire/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	beego.BConfig.WebConfig.EnableDocs = true
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/questionaire.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	beego.Run()
}
