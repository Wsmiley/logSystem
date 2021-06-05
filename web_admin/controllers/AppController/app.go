package AppController

import (
	"fmt"
	"strings"
	"web_admin/components"
	model "web_admin/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type AppController struct {
	beego.Controller
}

//项目列表
func (p *AppController) AppList() {
	logs.Debug("enter index controller")

	p.Layout = "layout/layout.html"
	appList, err := model.GetAllAPPInfo()
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("服务器繁忙")
		p.TplName = "app/error.html"
		logs.Warn("get app list failed, err:%v", err)
		return
	}
	logs.Debug("get app list succ, data:%v", appList)
	p.Data["applist"] = appList
	p.TplName = "app/index.html"
}

//创建项目
func (p *AppController) AppApply() {

	logs.Debug("enter index controller")
	p.Layout = "layout/layout.html"
	p.TplName = "app/apply.html"
}

func (p *AppController) AppCreate() {

	logs.Debug("enter index controller")
	appName := p.GetString("app_name")
	appType := p.GetString("app_type")
	developPath := p.GetString("develop_path")
	ipstr := p.GetString("iplist")

	p.Layout = "layout/layout.html"

	if len(appName) == 0 || len(appType) == 0 || len(developPath) == 0 || len(ipstr) == 0 {
		p.Data["Error"] = fmt.Sprintf("非法参数")
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}

	appInfo := &model.AppInfo{}
	appInfo.AppName = appName
	appInfo.AppType = appType
	appInfo.DevelopPath = developPath
	appInfo.IP = strings.Split(ipstr, ",")

	err := model.CreateApp(appInfo)
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("创建项目失败，数据库繁忙")
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}

	p.Redirect("/app/list", 302)
}

func (p *AppController) AppDelete() {

	logs.Debug("enter index controller")
	appID := p.GetString("id")
	appName := p.GetString("app_name")
	p.Layout = "layout/layout.html"

	//通过appName名字查IP表，得到项目部署在哪些服务器中
	iplist, err := model.GetIPInfoByName(appName)
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("获取项目ip失败，数据库繁忙")
		p.TplName = "log/error.html"
		logs.Warn("invalid parameter")
		return
	}
	err = model.DeleteApp(appID)
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("删除项目失败，数据库繁忙")
		p.TplName = "app/error.html"

		logs.Warn("invalid parameter")
		return
	}

	keyFormat := components.WebLogKey
	for _, ip := range iplist {
		key := fmt.Sprintf(keyFormat, ip)
		err = model.DeleteKeyToEtcd(key)
		if err != nil {
			logs.Warn("delete keyConf to etcd failed, err:%v", err)
			continue
		}
	}

	p.Redirect("/app/list", 302)
}
