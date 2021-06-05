package tailfile

import (
	"logAgent/common"

	"github.com/sirupsen/logrus"
)

//tailTask管理者

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask       //所有tailTask任务
	collectEntryList []common.CollectEntry      //所有配置项
	confChan         chan []common.CollectEntry //等待新配置的通道
}

var (
	ttMgr *tailTaskMgr
)

func Init(allConf []common.CollectEntry) (err error) {
	//allConf有若干个不不同的日志收集项目
	//每一个创建一个对应的tailObj
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		collectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry),
	}
	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Errorf("create tailobj for path:%s faild,err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt //把创建的这个tailTask任务登记，方便后续管理
		//启动goroutine去收集日志
		go tt.run()
	}
	go ttMgr.watch() //等新配置
	return
}

func (t *tailTaskMgr) watch() {
	for {
		//等待新配置
		newConf := <-t.confChan
		logrus.Infof("get new conf from etcd,conf:%v start manage tailTask", newConf)
		for _, conf := range newConf {
			//1. 原来存在的不用操作
			if t.isExist(conf) {
				continue
			}
			//2. 原来没有的，新创建一个tailTask
			tt := newTailTask(conf.Path, conf.Topic)
			err := tt.Init()
			if err != nil {
				logrus.Errorf("create tailobj for path:%s faild,err:%v", conf.Path, err)
				continue
			}
			logrus.Infof("create a tail task for path:%s success", conf.Path)
			ttMgr.tailTaskMap[tt.path] = tt //把创建的这个tailTask任务登记，方便后续管理
			//启动goroutine去收集日志
			go tt.run()

		}
		//3.原来有的现在没有的要tailTask停掉
		//找出tailTaskMap中存在，但newConf不存在的那些tailTask,把它们停掉
		for key, task := range t.tailTaskMap {
			var found bool
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				//这个tailTask需要结束
				logrus.Infof("the Task collect path:%s  need to stop", task.path)
				delete(t.tailTaskMap, key) //从tailTaskMap中删除
				task.cancel()
			}
		}
	}
}

//判断tailTaskMap中是否存在该收集项
func (t *tailTaskMgr) isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}

func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
