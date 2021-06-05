package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd faild,err:%v", err)
		return
	}
	defer cli.Close()
	//put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//str := `[{"path":"./logs/my.log","topic":"s4_log"},{"path":"./logs/my1.log","topic":"web_log"},{"path":"./logs/my2.log","topic":"s3_log"},{"path":"./logs/my3.log","topic":"s2_log"},{"path":"./logs/nignx.log","topic":"nignx"}]`
	//str := `[{"path":"./logs/my.log","topic":"s4_log"},{"path":"./logs/my1.log","topic":"web_log"},{"path":"./logs/my2.log","topic":"s3_log"},{"path":"./logs/my3.log","topic":"s2_log"}]`
	//str := `[{"path":"./logs/my.log","topic":"s4_log"},{"path":"./logs/my1.log","topic":"web_log"},{"path":"./logs/my2.log","topic":"s3_log"}]`
	//put(context,key,value)
	//_, err = cli.Put(ctx, "collect_log_127.0.0.1_conf", str)
	//cancel()
	if err != nil {
		fmt.Printf("put to etcd faild,err:%v\n", err)
		return
	}
	//get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "collect_log_127.0.0.4_conf")
	if err != nil {
		fmt.Printf("get from etcd faild,err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("key:%s value:%s\n", ev.Key, ev.Value)
	}
	cancel()
}
