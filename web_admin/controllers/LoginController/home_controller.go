package LoginController

import "fmt"

type HomeController struct {
	//beego.Controller
	BaseController
}

func (p *HomeController) Get() {
	fmt.Println("IsLogin:", p.IsLogin, p.Loginuser)
	if p.IsLogin != false {
		p.TplName = "layout/layout.html"
	} else {
		p.TplName = "login/login.html"
	}

}
