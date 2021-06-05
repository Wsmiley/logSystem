package components

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/logs/es"
)

var (
	BConfig   config.Configer = nil
	WebLogKey string
)

//初始化配置文件
func init() {
	var err error
	BConfig, err = config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println("config init error:", err)
		return
	}
	WebLogKey = BConfig.String("log::log_key")
}

func InitLogger() (err error) {
	if BConfig == nil {
		err = errors.New("beego config new failed!")
		return
	}
	maxlines, lerr := BConfig.Int64("log::maxlines")
	if lerr != nil {
		maxlines = 1000
	}

	logConf := make(map[string]interface{})
	logConf["filename"] = BConfig.String("log::log_path")
	level, _ := BConfig.Int("log::log_level")
	logConf["level"] = level
	logConf["maxlines"] = maxlines

	confStr, err := json.Marshal(logConf)
	_, err = json.Marshal(logConf)
	if err != nil {
		fmt.Println("logConf marshal failed,err:", err)
		return
	}

	err = logs.SetLogger(logs.AdapterFile, string(confStr))
	if err != nil {
		panic(err)
	}
	err = logs.SetLogger(logs.AdapterEs, `{"dsn":"http://localhost:9200/","level":7}`)
	if err != nil {
		panic(err)
	}
	logs.SetLogFuncCall(true)
	return
}
