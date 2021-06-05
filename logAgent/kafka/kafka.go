package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

type Message struct {
	Data  string
	Topic string
}

//初始化区全局kafka连接
func Init(address []string, ChanSize int64) (err error) {
	//1.生产者配置

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //ACK
	config.Producer.Partitioner = sarama.NewRandomPartitioner //分区
	config.Producer.Return.Successes = true                   //确认

	//2.连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka:producer closed,err:", err)
		return
	}
	//读日志文件和发送到kafka用msgChan做成异步操作,相当于缓冲
	msgChan = make(chan *sarama.ProducerMessage, ChanSize)
	go sendMsg()
	return
}

//从msgChan中读取msg，发送给kafka
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg faild,err:", err)
				return
			}
			logrus.Infof("send msg to kafka success,pid:%v offsset：%v", pid, offset)
		}
	}
}

//向外暴露msgChan
func ToMsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
