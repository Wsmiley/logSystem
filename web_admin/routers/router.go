package routers

import (
	"web_admin/controllers/AppController"
	"web_admin/controllers/LogController"
	"web_admin/controllers/LoginController"

	"github.com/astaxie/beego"
)

/*
/index为首页展示
/app/list为项目列表
/app/apply为创建项目
/app/create为创建项目后跳转的路由
/log/apply为创建日志
/log/list为为日志列表展示
/log/create为日志创建成功后跳转的路由
*/

func init() {

	// 登陆路由
	beego.Router("/home", &LoginController.HomeController{})

	beego.Router("/user/login", &LoginController.LoginController{}, "*:Login")
	beego.Router("/user/exit", &LoginController.ExitController{})

	beego.Router("/app/list", &AppController.AppController{}, "*:AppList")
	beego.Router("/app/apply", &AppController.AppController{}, "*:AppApply")
	beego.Router("/app/create", &AppController.AppController{}, "*:AppCreate")
	beego.Router("/app/:delete", &AppController.AppController{}, "*:AppDelete")

	beego.Router("/log/apply", &LogController.LogController{}, "*:LogApply")
	beego.Router("/log/ip/:id", &LogController.LogController{}, "*:LogIp")
	beego.Router("/log/list", &LogController.LogController{}, "*:LogList")
	beego.Router("/log/create", &LogController.LogController{}, "*:LogCreate")
	beego.Router("/log/:delete", &LogController.LogController{}, "*:LogDelete")
}
