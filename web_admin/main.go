package main

import (
	"strings"
	"web_admin/components"
	_ "web_admin/routers"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

var FilterUser = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("loginuser").(string)
	ok2 := strings.Contains(ctx.Request.RequestURI, "/user/login")
	if !ok && !ok2 {
		ctx.Redirect(302, "/user/login")
	}
}

func init_components() bool {
	err := components.InitLogger()
	if err != nil {
		logs.Warn("initDb failed, err :%v", err)
		return false
	}

	err = components.InitDb()
	if err != nil {
		logs.Warn("initDb failed, err:%v", err)
		return false
	}

	err = components.InitEtcd()

	if err != nil {
		logs.Warn("init etcd failed, err:%v", err)
		return false
	}
	return true

}

func main() {
	if init_components() == false {
		return
	}
	//注册过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)

	//开启session
	beego.BConfig.WebConfig.Session.SessionOn = true

	beego.Run()
}
