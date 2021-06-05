package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"logtransfer/kafka"
	"logtransfer/model"
	"time"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd faild,err:%v", err)
		return
	}
	return
}

//拉去日志配置项函数
func GetConf(key string) (collectEntryList []model.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key:%s faild,err:%v", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len:0 from etcd by key:%s", key)
	}
	keyValue := resp.Kvs[0]

	err = json.Unmarshal(keyValue.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("json.Unmarshal faild,err:%v", key, err)
		return
	}
	logrus.Debugf("load conf from etcd success,conf:%#v", collectEntryList)
	return
}

// 监控etcd中日志收集项配置变化
func WatchConf(key string) {
	for {
		watchCh := client.Watch(context.Background(), key)
		for wresp := range watchCh {
			logrus.Infof("get new conf from etcd!")
			for _, evt := range wresp.Events {
				fmt.Printf("type:%s key:%s value:%s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
				var newConf []model.CollectEntry
				if evt.Type == clientv3.EventTypeDelete {
					logrus.Warning("etcd delete the key!")
					kafka.SendNewConf(newConf) //没有接收即阻塞
					continue
				}
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Error("json unmarshal new conf failed,err:%v", err)
					continue
				}

				//有新配置后，应该告诉tailfile这个模块要启用新配置了
				kafka.SendNewConf(newConf) //没有接收即阻塞
			}
		}
	}
}
