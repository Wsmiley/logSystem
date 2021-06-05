package LoginController

import (
	"fmt"
	model "web_admin/models"

	"github.com/astaxie/beego"
	"github.com/beego/beego/v2/core/logs"
)

type LoginController struct {
	beego.Controller
}

func (p *LoginController) Get() {
	p.TplName = "login/login.html"
}

//admin   123456
func (p *LoginController) Login() {
	useNname := p.GetString("username")
	passWord := p.GetString("password")
	if useNname == "" || passWord == "" {
		logs.Info("input data is not valid")
		p.TplName = "login/login.html"
		return
	}
	id := model.QueryUserWithParam(useNname, passWord)
	fmt.Println("id:", id)
	if id > 0 {
		/*
			设置了session后会将数据处理设置到cookie，然后再浏览器进行网络请求的时候回自动带上cookie
			因为我们可以通过获取这个cookie来判断用户是谁，这里我们使用的是session的方式进行设置
		*/
		v := p.GetSession("loginuser")
		if v == nil {
			p.SetSession("loginuser", useNname)
			p.Data["num"] = 0

		} else {
			p.SetSession("asta", v.(int)+1)
			p.Data["num"] = v.(int)
		}
		p.Redirect("/home", 302)
	} else {
		p.Data["json"] = map[string]interface{}{"code": 0, "message": "登录失败"}
		p.ServeJSON()
	}
}
