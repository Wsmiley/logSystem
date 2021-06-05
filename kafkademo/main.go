package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func main() {
	//1.生产者配置

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //ACK
	config.Producer.Partitioner = sarama.NewRandomPartitioner //分区
	config.Producer.Return.Successes = true                   //确认

	//3.连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9095", "127.0.0.1:9094", "127.0.0.1:9093"}, config)
	if err != nil {
		fmt.Println("producer closed,err:", err)
	}
	defer client.Close()

	//2. 封装消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "CpuInfo"
	msg.Value = sarama.StringEncoder("2019.8.14Go")

	//4. 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg faild,err", err)
	}
	fmt.Printf("pid:%v offset：%v\n", pid, offset)

}
