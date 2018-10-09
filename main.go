package main

import (
	_ "LoveHome/routers"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
	"github.com/astaxie/beego/context"
)

func main() {
	beego.Run()
}
//访问静态资源
func ignoreStaticPath()  {
	//透明static
	beego.InsertFilter("/",beego.BeforeRouter,TransparentStatic)
	beego.InsertFilter("/*",beego.BeforeRouter,TransparentStatic)
}
func TransparentStatic(ctx *context.Context)  {
	//获取path 找不到方法
	orpath:=ctx.Request.URL.Path
	beego.Debug(" request url",orpath)
	if strings.Index(orpath,"api")>=0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter,ctx.Request,"static/html/"+ctx.Request.URL.Path)
}

