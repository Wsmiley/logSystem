package tailfile

import (
	"context"
	"fmt"
	"logAgent/kafka"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

func newTailTask(path, topic string) *tailTask {
	ctx, cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	return tt
}

func (t *tailTask) Init() (err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tObj, err = tail.TailFile(t.path, cfg)
	if err != nil {
		logrus.Error("tailfile create tailObj for path %s faild,err:%v\n", t.path, err)
		return
	}
	return
}

func (t *tailTask) run() (err error) {
	//读取日志发往kakfa
	logrus.Infof("collect for path:%s is running ", t.path)
	for {
		select {
		case <-t.ctx.Done(): //调用context中的cancel()方法就会收到信号
			logrus.Infof("path:%s is stopping", t.path)
			return
		case line, ok := <-t.tObj.Lines:
			if !ok {
				logrus.Warn("tail file close reopen path %s\n", t.path)
				time.Sleep(time.Second)
				continue
			}
			//空行则略过
			if len(strings.Trim(line.Text, "\r")) == 0 {
				continue
			}
			//利用通道，将同步的代码改为异步
			//把读出来的一行日志包装成kafka里面的msg类型，丢到通道中
			fmt.Println("log:", line.Text)
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			kafka.ToMsgChan(msg)
		}

	}
	return
}
