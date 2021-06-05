package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	rand1 "math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
//模拟生成Nignx日志样例
{"@timestamp":"31/Mar/2020:12:10:23 +0800","host": "192.168.20.2","clientip": "192.168.20.88","size": 82896,"responsetime": 0.006,"upstreamtime": "0.006","upstreamhost": "192.168.20.3:80","http_host": "192.168.20.2","url": "/img/header-background.png","xff": "-","referer": "http://192.168.20.2/","agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36","status": "200"}
*/

//资源结构体
type urlJson struct {
	Timestamp    string `json:"@timestamp"`      //请求时间
	Host         string `json:"host"`            //服务器地址
	Remote_addr  string `json:"remote_addr"`     //客户端的IP地址。
	Request      string `json:"request"`         //请求与http协议
	Status       string `json:"status"`          //请求状态
	Bytes        string `json:"body_bytes_sent"` //传送页面的字节数
	Http_refer   string `json:"http_refer"`
	UserAgent    string `json:"http_user_agent"` //客户端浏览器相关信息
	Request_time string `json:"request_time"`    //请求处理时间，单位为秒，精度毫秒
}

//常用的userAgent收集
var userAgentList = []string{

	//Android平台原生浏览器
	"Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166  Safari/535.19",
	"Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; GT-I9300 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
	"Mozilla/5.0 (Linux; U; Android 2.2; en-gb; GT-P1000 Build/FROYO) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",

	//Firefox火狐
	"Mozilla/5.0 (Android; Mobile; rv:14.0) Gecko/14.0 Firefox/14.0",
	"Mozilla/5.0 (Android; Tablet; rv:14.0) Gecko/14.0 Firefox/14.0",
	"Mozilla/5.0 (Windows NT 6.2; WOW64; rv:21.0) Gecko/20100101 Firefox/21.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:21.0) Gecko/20100101 Firefox/21.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:21.0) Gecko/20130331 Firefox/21.0",

	//Google chrome
	"Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19",
	"Mozilla/5.0 (Linux; Android 4.1.2; Nexus 7 Build/JZ054K) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166 Safari/535.19",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.11 (KHTML, like Gecko) Ubuntu/11.10 Chromium/27.0.1453.93 Chrome/27.0.1453.93 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.94 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 6_1_4 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) CriOS/27.0.1453.10 Mobile/10B350 Safari/8536.25",

	//Internet Explore
	"Mozilla/5.0 (compatible; WOW64; MSIE 10.0; Windows NT 6.2)",      //IE10
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)", //IE9
	"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)", //IE8
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)",              //IE7
	"Mozilla/4.0 (Windows; MSIE 6.0; Windows NT 5.2)",                 //IE6

	//Opera
	"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.9.168 Version/11.52", //Mac
	"Opera/9.80 (Windows NT 6.1; WOW64; U; en) Presto/2.10.229 Version/11.62",           //Windows

	//Safari
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_6; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27",      //Mac
	"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27",               //windows
	"Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",          //iPad
	"Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3", //iPhone

	//iOS
	"Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",         //iPad
	"Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.", //iPhone
	"Mozilla/5.0 (iPod; U; CPU like Mac OS X; en) AppleWebKit/420.1 (KHTML, like Gecko) Version/3.0 Mobile/3A101a Safari/419.3",             //iPod

	//Windows Phone
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows Phone OS 7.0; Trident/3.1; IEMobile/7.0; LG; GW910)",                   //windows Phone 7
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0; SAMSUNG; SGH-i917)",           // Windows Phone7.5
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows Phone 8.0; Trident/6.0; IEMobile/10.0; ARM; Touch; NOKIA; Lumia 920)", //windows phone 8
}
var url_path_list = []string{
	"login.php",
	"view.php",
	"list.php",
	"upload.php",
	"admin/login.php",
	"edit.php",
	"index.html",
}

var ip_slice_list = []uint8{10, 29, 30, 46, 55, 63, 72, 87, 98, 132, 156, 124, 167, 143, 187, 168, 190, 201, 202, 214, 215, 222}

var status = []string{"200", "301", "404", "500"}

var http_refer = []string{"http://www.baidu.com/s?wd={query}", "http://www.google.cn/search?q={query}", "http://www.sogou.com/web?query={query}", "http://one.cn.yahoo.com/s?p={query}", "http://cn.bing.com/search?q={query}"}

var search_keyword = []string{"spark", "hadoop", "hive", "spark mlib", "spark sql"}

//使用时间作为随机种子，在不同时间下随机出来的结果是不一样的
//不设置随机种子的情况下，就会出现伪随机数

func randIP() string {
	var str string
	for i := 0; i < 4; i++ {
		result, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ip_slice_list))))
		b := strconv.Itoa(int(ip_slice_list[result.Uint64()]))
		if i < 3 {
			str += b + "."
		} else {
			str += b
		}

	}
	return str
}

func randUrlPath() string {
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(len(url_path_list))))
	return url_path_list[result.Uint64()]
}

func randStatus() string {
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(len(status))))
	if rand1.Float32() > 0.2 {
		return status[0]
	}
	return status[result.Uint64()]
}

func randBytes() string {
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(520)))
	return result.String()
}

func randUserAgent() string {
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userAgentList))))
	return userAgentList[result.Uint64()]
}

func randRequestTime() string {
	timens := int64(time.Now().Nanosecond())
	rand1.Seed(timens)
	return strconv.FormatFloat(float64(0.1*rand1.Float32()), 'f', 2, 64)
}

func randRefer() string {
	replace := "{query}"
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(len(search_keyword))))
	if rand1.Float32() > 0.2 {
		str := http_refer[0]
		str = strings.Replace(str, replace, search_keyword[result.Uint64()], 1)
		return str
	}
	result1, _ := rand.Int(rand.Reader, big.NewInt(int64(len(http_refer))))
	str := http_refer[result1.Uint64()]
	str = strings.Replace(str, replace, search_keyword[result.Uint64()], 1)
	return str

}

func WriteWithIoutil(name string, contents string) {
	file, _ := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	defer file.Close()
	// 获取bufio.Writer实例
	writer := bufio.NewWriter(file)
	// 写入字符串
	count, _ := file.WriteString(contents)
	fmt.Printf("wrote %d byte\n", count)
	// 清空缓存 确保写入磁盘
	writer.Flush()
}

func makelog() {
	log := urlJson{
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		Host:         "127.0.0.1",
		Remote_addr:  randIP(),
		Request:      "Get /" + randUrlPath() + " HTTP/1.1",
		Status:       randStatus(),
		Bytes:        randBytes(),
		Http_refer:   randRefer(),
		UserAgent:    randUserAgent(),
		Request_time: randRequestTime(),
	}
	jsons, errs := json.Marshal(log) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	jsondata := string(jsons)
	jsondata = jsondata + "\n"
	WriteWithIoutil("D:/project/logSystem/logAgent/logs/test.log", jsondata)
}

func main() {
	for i := 0; i < 100; i++ {
		makelog()
		time.Sleep(100 * time.Millisecond)
	}
}
