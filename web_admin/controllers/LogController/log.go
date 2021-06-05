package LogController

import (
	"fmt"
	"strconv"
	"web_admin/components"
	model "web_admin/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type LogController struct {
	beego.Controller
}

type JsonData struct {
	Code int
	Msg  string
	Data []string
}

func (p *LogController) LogList() {

	logs.Debug("enter index controller")

	p.Layout = "layout/layout.html"
	logList, err := model.GetAllLogInfo()
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("服务器繁忙")
		p.TplName = "app/error.html"

		logs.Warn("get app list failed, err:%v", err)
		return
	}

	logs.Debug("get app list succ, data:%v", logList)
	p.Data["loglist"] = logList

	p.TplName = "log/index.html"
}

func (p *LogController) LogApply() {

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
	p.Data["appList"] = appList

	p.TplName = "log/apply.html"
}

func (p *LogController) LogCreate() {

	logs.Debug("enter index controller")
	appName := p.GetString("app_name")
	logPath := p.GetString("log_path")
	topic := p.GetString("topic")
	ip_test := p.GetString("app_ip")
	p.Layout = "layout/layout.html"
	if len(appName) == 0 || len(logPath) == 0 || len(topic) == 0 {
		p.Data["Error"] = fmt.Sprintf("非法参数")
		p.TplName = "log/error.html"

		logs.Warn("invalid parameter")
		return
	}
	logInfo := &model.LogInfo{}
	logInfo.AppName = appName
	logInfo.LogPath = logPath
	logInfo.Topic = topic
	logInfo.Ip = ip_test
	err := model.CreateLog(logInfo)
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("创建项目失败，数据库繁忙")
		p.TplName = "log/error.html"

		logs.Warn("invalid parameter")
		return
	}

	//通过appName名字查IP表，得到项目部署在哪些服务器中
	// iplist, err := model.GetIPInfoByName(appName)
	// if err != nil {
	// 	p.Data["Error"] = fmt.Sprintf("获取项目ip失败，数据库繁忙")
	// 	p.TplName = "log/error.html"

	// 	logs.Warn("invalid parameter")
	// 	return
	// }

	if err != nil {
		p.Data["Error"] = fmt.Sprintf("获取项目下的Log失败，数据库繁忙")
		p.TplName = "log/error.html"

		logs.Warn("invalid parameter,err:%v", err)
		return
	}
	keyFormat := components.WebLogKey
	key := fmt.Sprintf(keyFormat, ip_test)
	err = model.SetLogConfToEtcd(key, ip_test, "create")
	if err != nil {
		logs.Warn("Set log conf to etcd failed, err:%v", err)
	}
	// for _, ip := range iplist {
	// 	key := fmt.Sprintf(keyFormat, ip)
	// 	//通过和key获取先获取到VALUE，再添加Value，避免覆盖

	// 	err = model.SetLogConfToEtcd(key, logINFO)
	// 	if err != nil {
	// 		logs.Warn("Set log conf to etcd failed, err:%v", err)
	// 		continue
	// 	}
	// }
	p.Redirect("/log/list", 302)
}

func (p *LogController) LogDelete() {
	logs.Debug("enter index controller")
	logID := p.GetString("id")
	appName := p.GetString("app_name")
	logPath := p.GetString("log_path")
	topic := p.GetString("topic")
	ip := p.GetString("ip")

	p.Layout = "layout/layout.html"
	logInfo := &model.LogInfo{}
	logInfo.LogId, _ = strconv.Atoi(logID)
	logInfo.LogPath = logPath
	logInfo.AppName = appName
	logInfo.Topic = topic
	logInfo.Ip = ip
	err := model.DeleteLog(logInfo)
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("删除失败，数据库繁忙")
		p.TplName = "log/error.html"

		logs.Warn("invalid parameter")
		return
	}

	keyFormat := components.WebLogKey
	key := fmt.Sprintf(keyFormat, ip)
	err = model.SetLogConfToEtcd(key, ip, "delete")
	if err != nil {
		logs.Warn("Set log conf to etcd failed, err:%v", err)
	}
	p.Redirect("/log/list", 302)
}

func (p *LogController) LogIp() {

	logs.Debug("enter index controller")
	AppID := p.GetString("appid")
	fmt.Println(AppID)
	appid, err := strconv.Atoi(AppID)
	IpList, err := model.GetIPInfoById(appid)
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("服务器繁忙")
		p.TplName = "log/error.html"
		logs.Warn("get ip list failed, err:%v", err)
		return
	}
	logs.Debug("get ip list succ, data:%v", IpList)
	p.Data["json"] = JsonData{200, "获取成功", IpList}
	p.ServeJSON()
}
