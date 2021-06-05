package models

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
)

type AppInfo struct {
	AppId       int      `db:"app_id"`
	AppName     string   `db:"app_name"`
	AppType     string   `db:"app_type"`
	CreateTime  string   `db:"create_time"`
	DevelopPath string   `db:"develop_path"`
	IP          []string `db:"ip"`
}

var (
	Db *sqlx.DB
)

//初始化数据库
func InitDb(db *sqlx.DB) {
	Db = db
}

func GetAllAPPInfo() (appList []AppInfo, err error) {
	err = Db.Select(&appList, "select app_id,app_name,app_type,create_time,develop_path from tbl_app_info")
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}

	var ipMap map[int][]string
	ipMap = make(map[int][]string)
	for _, v := range appList {
		var ipList []string
		err = Db.Select(&ipList, "select ip from tbl_app_ip where tbl_app_ip.app_id=?", v.AppId)
		ipMap[v.AppId] = ipList
	}
	for _, v := range appList {
		v.IP = ipMap[v.AppId]
	}
	return
}

func GetIPInfoById(appId int) (iplist []string, err error) {
	err = Db.Select(&iplist, "select ip from tbl_app_ip where app_id=?", appId)
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func GetIPInfoByName(appName string) (iplist []string, err error) {

	var appId []int
	err = Db.Select(&appId, "select app_id from tbl_app_info where app_name=?", appName)
	if err != nil || len(appId) == 0 {
		logs.Warn("select app_id failed, Db.Exec error:%v", err)
		return
	}

	err = Db.Select(&iplist, "select ip from tbl_app_ip where app_id=?", appId[0])
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func CreateApp(info *AppInfo) (err error) {

	conn, err := Db.Begin()
	if err != nil {
		logs.Warn("CreateApp failed, Db.Begin error:%v", err)
		return
	}
	//失败则回滚，成功则提交
	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}

		conn.Commit()
	}()

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	r, err := conn.Exec("insert into tbl_app_info(app_name, app_type, develop_path,create_time)values(?, ?, ?, ?)",
		info.AppName, info.AppType, info.DevelopPath, timeStr)

	if err != nil {
		logs.Warn("CreateApp failed, Db.Exec error:%v", err)
		return
	}

	id, err := r.LastInsertId()
	if err != nil {
		logs.Warn("CreateApp failed, Db.LastInsertId error:%v", err)
		return
	}

	//插入IP
	for _, ip := range info.IP {
		_, err = conn.Exec("insert into tbl_app_ip(app_id, ip,create_time)values(?,?,?)", id, ip, timeStr)
		if err != nil {
			logs.Warn("CreateApp failed, conn.Exec ip error:%v", err)
			return
		}
	}
	return
}

func DeleteApp(appID string) (err error) {
	conn, err := Db.Begin()
	if err != nil {
		logs.Warn("DeleteAPP failed, Db.Begin error:%v", err)
		return
	}

	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}

		conn.Commit()
	}()
	//删除app_info表的数据
	sqlStr := "DELETE FROM tbl_app_info WHERE app_id = ?"
	result, err := Db.Exec(sqlStr, appID)
	if err != nil {
		logs.Warn("exec failed, err:%v\n", err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		logs.Warn("get affected failed, err:%v\n", err)
		return
	}
	//删除log_info表的数据
	sqlStr1 := "DELETE FROM tbl_log_info WHERE app_id = ?"
	result1, err := Db.Exec(sqlStr1, appID)
	if err != nil {
		logs.Warn("exec failed, err:%v\n", err)
		return
	}
	_, err = result1.RowsAffected()
	if err != nil {
		logs.Warn("get affected failed, err:%v\n", err)
		return
	}
	//删除app_ip表的数据
	sqlStr2 := "DELETE FROM tbl_app_ip WHERE app_id = ?"
	result2, err := Db.Exec(sqlStr2, appID)
	if err != nil {
		logs.Warn("exec failed, err:%v\n", err)
		return
	}
	_, err = result2.RowsAffected()
	if err != nil {
		logs.Warn("get affected failed, err:%v\n", err)
		return
	}
	return
}
