package main

import (
	"fmt"
	"logAgent/common"
	"logAgent/etcd"
	"logAgent/kafka"
	"logAgent/sysinfo"
	"logAgent/tailfile"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

//日志收集的客户端
//收集指定目录下的日志文件，发送到kafka中

type Config struct {
	KafaConfig    `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
	InfluxConfig  `ini:"influx"`
}

type KafaConfig struct {
	Address  string `ini:"address"`
	ChanSize int64  `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

type InfluxConfig struct {
	Address  string `ini:"address"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	Database string `ini:"database"`
}

func run() {
	select {}
}

func main() {
	//-1.获取本机ip，为Etcd获取配置文件
	ip, err := common.GetOutboundIP()
	if err != nil {
		logrus.Errorf("get ip faild,err:%v", err)
	}
	ip = "127.0.0.1"
	//0.初始化,读配置文件
	var configOBJ = new(Config)
	err = ini.MapTo(configOBJ, "./conf/config.ini")
	if err != nil {
		logrus.Errorf("load config faild,err：%v", err)
		return
	}
	fmt.Printf("%#v\n", configOBJ)
	err = kafka.Init([]string{configOBJ.KafaConfig.Address}, configOBJ.KafaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka faild,err：%v", err)
		return
	}
	logrus.Info("init kafka success!")

	//1.根据配置中的日志路径使用tail收集数据
	//初始化etcd
	err = etcd.Init([]string{configOBJ.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd faild,err:%v", err)
		return
	}
	//从etcd中拉去日志配置项
	collectKey := fmt.Sprintf(configOBJ.EtcdConfig.CollectKey, ip)
	allConf, err := etcd.GetConf(collectKey)
	fmt.Printf("allConf:%s\n", allConf)
	if err != nil {
		logrus.Errorf("get Conf from etcd faild,err:%v", err)
		return
	}
	//监控etcd中，configObj.EtcdConfig.Colloctkey对应值的变化
	go etcd.WatchConf(collectKey)
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Errorf("init tail faild,err:%v", err)
		return
	}
	logrus.Infof("init tail success!")
	//sysinfo collect
	addr := configOBJ.InfluxConfig.Address
	username := configOBJ.InfluxConfig.Username
	password := configOBJ.InfluxConfig.Password
	database := configOBJ.InfluxConfig.Database
	err = sysinfo.InitconnInflux(addr, username, password, database)
	if err != nil {
		logrus.Errorf("init influxdb faild,err:%v", err)
		return
	}
	go sysinfo.Run(time.Second)
	run()
}
