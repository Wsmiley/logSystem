package model

type Config struct {
	KafkaConf `ini:"kafka"`
	ESConf    `ini:"es"`
	EtcdConf  `ini:"etcd"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type ESConf struct {
	Address string `ini:"address"`
	MaxSize int    `ini:"max_chan_size"`
	GoNum   int    `ini:"gorountine_num"`
}

type EtcdConf struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

//日志配置项结构体
type CollectEntry struct {
	Path  string `json:"path"`  //日志文件路径
	Topic string `json:"topic"` //日志数据发往那个topic
}
