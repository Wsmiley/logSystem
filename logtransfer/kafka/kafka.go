package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"logtransfer/es"
	"logtransfer/model"
	"time"

	"github.com/Shopify/sarama"

	"github.com/sirupsen/logrus"
)

var (
	consumer sarama.Consumer
	kkMgr    *kafkaTaskMgr
)

type kafkaTask struct {
	topic  string
	ctx    context.Context
	cancel context.CancelFunc
}

type kafkaTaskMgr struct {
	kafkaTaskMap     map[string]*kafkaTask     //所有tailTask任务
	collectEntryList []model.CollectEntry      //所有配置项
	confChan         chan []model.CollectEntry //等待新配置的通道
}

type logJsonData struct {
	TimeStamp     string `json:"timestamp"`
	Msg           string `json:"msg"`
	TimeStampLog  string `json:"@timestamp"`
	Host          string `json:"host"`
	RemoteAddr    string `json:"remote_addr"`
	Request       string `json:"request"`
	Status        string `json:"status"`
	BodyBytesSent string `json:"body_bytes_sent"`
	HTTPRefer     string `json:"http_refer"`
	HTTPUserAgent string `json:"http_user_agent"`
	RequestTime   string `json:"request_time"`
}

func Init(addr []string) (err error) {
	consumer, err = sarama.NewConsumer(addr, nil)
	if err != nil {
		fmt.Printf("faild to start consumer,err:%v\n", err)
		return
	}
	return
}

func newKafkaTask(topic string) *kafkaTask {
	ctx, cancel := context.WithCancel(context.Background())
	kk := &kafkaTask{
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	return kk
}

func RecvKafka(allConf []model.CollectEntry) (err error) {
	kkMgr = &kafkaTaskMgr{
		kafkaTaskMap:     make(map[string]*kafkaTask, 20),
		collectEntryList: allConf,
		confChan:         make(chan []model.CollectEntry),
	}
	for _, conf := range allConf {
		k := newKafkaTask(conf.Topic)
		kkMgr.kafkaTaskMap[k.topic] = k //把创建的这个tailTask任务登记，方便后续管理
		k.run()
	}

	go kkMgr.watch() //等新配置

	return
}

func (k *kafkaTask) run() (err error) {

	//拿到指导topic下面所有分区列表
	partitionList, err1 := consumer.Partitions(k.topic) //根据topic取到所有分区
	if err1 != nil {
		fmt.Printf("faild to get list of Partitions,err:%v\n", err)
		return
	}

	for partition := range partitionList { //遍历所有分区
		//针对每个分区创建一个对应的消费者
		var pc sarama.PartitionConsumer
		pc, err1 = consumer.ConsumePartition(k.topic, int32(partition), sarama.OffsetNewest)
		if err1 != nil {
			fmt.Printf("faild to start consumer for partition %d ,err:%v\n", partition, err)
			return
		}
		go k.kafkaRun(pc)
		// 	go func(sarama.PartitionConsumer) {
		// 		logrus.Infof("collect for topic:%s is running ", k.topic)
		// 		for msg := range pc.Messages() {
		// 			//logDataChan-<msg //将同步流程异步化，将取出日志数据先放到channel中
		// 			//fmt.Println(msg.Topic, string(msg.Value))
		// 			var m1 map[string]interface{}
		// 			err = json.Unmarshal(msg.Value, &m1)
		// 			if err != nil {
		// 				fmt.Printf("Unmarshal msg faild,err:%v\n", err)
		// 				continue
		// 			}
		// 			es.PutLogData(m1)
		// 		}
		// 		defer pc.AsyncClose()
		// 	}(pc)
		// }
	}
	return
}

func (k *kafkaTask) kafkaRun(pc sarama.PartitionConsumer) {
	select {
	case <-k.ctx.Done():
		logrus.Infof("topic:%s is stopping", k.topic)
		return
	default:
		logrus.Infof("collect for topic:%s is running ", k.topic)
		for msg1 := range pc.Messages() {
			var m1 map[string]string

			err := json.Unmarshal(msg1.Value, &m1)
			if err != nil {
				fmt.Printf("Unmarshal msg faild,err:%v\n", err)
				idx := logJsonData{
					TimeStamp: time.Now().Format(time.RFC3339),
					Msg:       string(msg1.Value),
				}
				body, err := json.Marshal(idx)
				if err != nil {
					fmt.Printf("marshal msg faild,err:%v\n", err)
				}
				es.PutLogData(string(body), msg1.Topic)
				continue
			} else {
				var idx logJsonData
				idx.TimeStamp = time.Now().Format(time.RFC3339)
				for k, v := range m1 {
					if k == "@timestamp" {
						idx.TimeStampLog = v
					}
					if k == "host" {
						idx.Host = v
					}
					if k == "remote_addr" {
						idx.RemoteAddr = v
					}
					if k == "request" {
						idx.Request = v
					}
					if k == "status" {
						idx.Status = v
					}
					if k == "body_bytes_sent" {
						idx.BodyBytesSent = v
					}
					if k == "http_refer" {
						idx.HTTPRefer = v
					}
					if k == "http_user_agent" {
						idx.HTTPUserAgent = v
					}
					if k == "request_time" {
						idx.RequestTime = v
					}
				}
				body, err := json.Marshal(idx)
				if err != nil {
					fmt.Printf("marshal msg faild,err:%v\n", err)
				}
				es.PutLogData(string(body), msg1.Topic)
			}
		}
	}

}

func (k *kafkaTaskMgr) watch() {
	for {
		//等待新配置
		newConf := <-k.confChan
		logrus.Infof("get new conf from etcd,conf:%v start manage kafkaTask", newConf)
		for _, conf := range newConf {
			//1. 原来存在的不用操作
			if k.isExist(conf) {
				continue
			}
			//2. 原来没有的，新创建一个kafkaTask
			kk := newKafkaTask(conf.Topic)
			logrus.Infof("create a kafka collection task for topic:%s success", conf.Topic)
			kkMgr.kafkaTaskMap[kk.topic] = kk //把创建的这个tailTask任务登记，方便后续管理

			//启动goroutine去收集日志
			kk.run()

		}
		//3.原来有的现在没有的要kafkaTask停掉
		//找出kafkaTaskMgr中存在，但newConf不存在的那些kafkaTask,把它们停掉
		for key, task := range k.kafkaTaskMap {
			var found bool
			for _, conf := range newConf {
				if key == conf.Topic {
					found = true
					break
				}
			}
			if !found {
				//这个kafkaTask需要结束
				logrus.Infof("the Task collect topic:%s  need to stop", task.topic)
				delete(k.kafkaTaskMap, key) //从kafkaTaskMap中删除
				task.cancel()
			}
		}
	}
}

//判断kafkaTaskMap中是否存在该收集项
func (k *kafkaTaskMgr) isExist(conf model.CollectEntry) bool {
	_, ok := k.kafkaTaskMap[conf.Topic]
	return ok
}

func SendNewConf(newConf []model.CollectEntry) {
	kkMgr.confChan <- newConf
}
