package main

import (
	_ "github.com/udistrital/titan_api_mid/routers"
	"github.com/udistrital/utils_oas/apiStatusLib"
  	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego"
	//"github.com/udistrital/auditoria"

)


func main() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
	AllowHeaders: []string{"Origin", "x-requested-with",
	"content-type",
	"accept",
	"origin",
	"authorization",
	"x-csrftoken"},
	ExposeHeaders: []string{"Content-Length"},
	AllowCredentials: true,
	}))
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	//Hola cambio
	//InitInterceptor();
	//auditoria.InitMiddleware();
	apistatus.Init()
	beego.Run()

}
