package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	etcdclient "go.etcd.io/etcd/clientv3"
)

var (
	etcdClient *etcdclient.Client
)

type LogInfo struct {
	AppId      int    `db:"app_id"`
	AppName    string `db:"app_name"`
	LogId      int    `db:"log_id"`
	CreateTime string `db:"create_time"`
	LogPath    string `db:"log_path"`
	Topic      string `db:"topic"`
	Ip         string `db:"ip"`
}

// type EtcdLogConf struct {
// 	Path  string `json:"log_path"`
// 	Topic string `json:"topic"`
// }
type EtcdLogConf struct {
	Path  string `db:"log_path"`
	Topic string `db:"topic"`
}

//初始化ETCD
func InitEtcd(client *etcdclient.Client) {
	etcdClient = client
}

func GetAllLogInfo() (loglist []LogInfo, err error) {
	err = Db.Select(&loglist,
		"select a.app_id, b.app_name, a.create_time, a.topic, a.log_id, a.log_path,a.ip from tbl_log_info a, tbl_app_info b where a.app_id=b.app_id")
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func GetLogInfo(appName string) (loglist []LogInfo, err error) {
	str := "select b.app_name, a.topic,a.log_path from tbl_log_info a, tbl_app_info b where a.app_id=b.app_id and a.app_name=?"
	err = Db.Select(&loglist, str, appName)
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func GetLogIPInfo(Ip string) (loglist []EtcdLogConf, err error) {
	str := "select log_path,topic from tbl_log_info where ip=?"
	err = Db.Select(&loglist, str, Ip)
	if err != nil {
		logs.Warn("Get All IP Info failed, err:%v", err)
		return
	}
	return
}

func CreateLog(info *LogInfo) (err error) {

	conn, err := Db.Begin()
	if err != nil {
		logs.Warn("CreateLog failed, Db.Begin error:%v", err)
		return
	}

	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}

		conn.Commit()
	}()

	var appId []int

	err = Db.Select(&appId, "select app_id from tbl_app_info where app_name=?", info.AppName)
	if err != nil || len(appId) == 0 {
		logs.Warn("select app_id failed, Db.Exec error:%v", err)
		return
	}
	var logInfo []LogInfo
	err = Db.Select(&logInfo, "select app_id from tbl_log_info where app_name=?", info.AppName)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	info.AppId = appId[0]
	r, err := conn.Exec("insert into tbl_log_info(app_id, log_path, topic, app_name, create_time,ip)values(?, ?, ?,?,?,?)",
		info.AppId, info.LogPath, info.Topic, info.AppName, timeStr, info.Ip)

	if err != nil {
		logs.Warn("CreateLog failed, Db.Exec error:%v", err)
		return
	}

	_, err = r.LastInsertId()
	if err != nil {
		logs.Warn("CreateLog failed, Db.LastInsertId error:%v", err)
		return
	}
	return
}

func DeleteLog(info *LogInfo) (err error) {
	conn, err := Db.Begin()
	if err != nil {
		logs.Warn("DeletelOG failed, Db.Begin error:%v", err)
		return
	}

	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}

		conn.Commit()
	}()
	sqlStr := "DELETE FROM tbl_log_info WHERE log_id = ?"
	result, err := Db.Exec(sqlStr, info.LogId)
	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		fmt.Printf("get affected failed, err:%v\n", err)
		return
	}
	return
}

//collect_log_127.0.0.1_conf
func SetLogConfToEtcd(etcdKey string, ip string, choose string) (err error) {
	loglist, err := GetLogIPInfo(ip)
	data, err := json.Marshal(loglist)
	if err != nil {
		logs.Warn("marshal failed, err:%v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	if choose == "delete" {
		etcdClient.Delete(ctx, etcdKey)
	}
	_, err = etcdClient.Put(ctx, etcdKey, string(data))
	cancel()
	if err != nil {
		logs.Warn("Put failed, err:%v", err)
		return
	}

	logs.Debug("put etcd succ, data:%v", string(data))
	return
}

func DeleteKeyToEtcd(etcdKey string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = etcdClient.Delete(ctx, etcdKey)
	cancel()
	if err != nil {
		logs.Warn("delete failed, err:%v", err)
		return
	}
	return
}
