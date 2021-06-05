package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

// Elasticsearch demo

type logJsonData struct {
	TimeStamp string `json:"timestamp"`
	Msg       string `json:"msg"`
}

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("connect to es success")
	str1 := "{\"@timestamp\":\"2021-05-28 16:41:28\",\"host\":\"127.0.0.1\",\"remote_addr\":\"214.72.222.72\",\"request\":\"Get /view.php HTTP/1.1\",\"status\":\"200\",\"body_bytes_sent\":\"151\",\"http_refer\":\"http://www.baidu.com/s?wd=spark\",\"http_user_agent\":\"Mozilla/5.0 (Android; Tablet; rv:14.0) Gecko/14.0 Firefox/14.0\",\"request_time\":\"0.01\"}"
	put1, err := client.Index().
		Index("test").
		BodyJson(str1).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
