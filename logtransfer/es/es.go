package es

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

//将日志数据谢瑞Elasticsearch
type ESClient struct {
	client      *elastic.Client
	indexChan   chan string
	logDataChan chan string
}

var (
	esClient *ESClient
)

func Init(addr string, gorountineNum int, maxSize int) (err error) {

	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://"+addr))
	esClient = &ESClient{
		client:      client,
		indexChan:   make(chan string, maxSize),
		logDataChan: make(chan string, maxSize),
	}
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("connect to es success")
	for i := 0; i < gorountineNum; i++ {
		go sendToES()
	}

	return
}

func sendToES() {
	for msg := range esClient.logDataChan {
		put1, err := esClient.client.Index().
			Index(<-esClient.indexChan).
			BodyString(msg).
			Do(context.Background())
		if err != nil {
			// Handle error
			fmt.Printf("send to es failed, err: %v\n", err)
			continue
		}
		fmt.Printf("Indexed user %s to index %s,type %s\n", put1.Id, put1.Index, put1.Type)
	}

}

func PutLogData(msg string, topic string) {
	esClient.indexChan <- topic
	esClient.logDataChan <- msg
}
