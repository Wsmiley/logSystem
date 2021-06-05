package LoginController

import (
	"fmt"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	IsLogin   bool
	Loginuser interface{}
}

//判断是否登录
func (p *BaseController) Prepare() {
	loginuser := p.GetSession("loginuser")
	fmt.Println("loginuser---->", loginuser)
	if loginuser != nil {
		p.IsLogin = true
		p.Loginuser = loginuser
	} else {
		p.IsLogin = false
	}
	p.Data["IsLogin"] = p.IsLogin
}
