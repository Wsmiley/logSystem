package components

import (
	"fmt"
	"time"
	model "web_admin/models"

	"go.etcd.io/etcd/clientv3"
)

func InitEtcd() (err error) {
	cli, err := clientv3.New(clientv3.Config{
		//Endpoints: []string{BConfig.String("etcd::etcd1"), BConfig.String("etcd::etcd2"),
		//BConfig.String("etcd::etcd3")},
		Endpoints:   []string{BConfig.String("etcd::etcd1")},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect etcd failed, err:", err)
		return
	}

	model.InitEtcd(cli)
	return
}
