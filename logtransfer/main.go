package main

//logtansfer
//从kafka消费数据

import (
	"fmt"
	"logtransfer/es"
	"logtransfer/etcd"
	"logtransfer/kafka"
	"logtransfer/model"

	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func main() {
	//1.配置文件
	var cfg = new(model.Config)
	err := ini.MapTo(cfg, "./config/logtransfer.ini")
	if err != nil {
		logrus.Errorf("load config failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("load config success")

	//2.连接Es
	err = es.Init(cfg.ESConf.Address, cfg.ESConf.GoNum, cfg.ESConf.MaxSize)
	if err != nil {
		logrus.Errorf("Init es failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("init es success")

	//3.连接kafka
	err = kafka.Init([]string{cfg.KafkaConf.Address})
	if err != nil {
		logrus.Errorf("connect to kafka failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("Init kafka  success")

	//4.etcd
	err = etcd.Init([]string{cfg.EtcdConf.Address})
	if err != nil {
		logrus.Errorf("init etcd faild,err:%v", err)
		return
	}

	ip := "127.0.0.1"
	collectKey := fmt.Sprintf(cfg.EtcdConf.CollectKey, ip)
	allConf, err := etcd.GetConf(collectKey)
	if err != nil {
		logrus.Errorf("get Conf from etcd faild,err:%v", err)
		return
	}
	fmt.Printf("allConf:%s\n", allConf)
	err = kafka.RecvKafka(allConf)

	if err != nil {
		logrus.Errorf("kafka allconf failed,err%:v", err)
		return
	}
	//监控etcd中，configObj.EtcdConfig.Colloctkey对应值的变化
	go etcd.WatchConf(collectKey)
	select {}
}
