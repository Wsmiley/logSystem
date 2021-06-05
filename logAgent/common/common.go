package common

import (
	"fmt"
	"net"
	"strings"
)

const (
	CanNotGetIp = "get ip faild"
)

//日志配置项结构体
type CollectEntry struct {
	Path  string `json:"path"`  //日志文件路径
	Topic string `json:"topic"` //日志数据发往那个topic
}

// Get preferred outbound ip of this machine
func GetOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}
