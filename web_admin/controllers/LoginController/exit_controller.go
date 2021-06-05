package LoginController

type ExitController struct {
	BaseController
}

func (p *ExitController) Get() {
	//清除该用户登录状态的数据
	p.DelSession("loginuser")
	p.TplName = "login/login.html"
}
