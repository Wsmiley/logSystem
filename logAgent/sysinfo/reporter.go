package sysinfo

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var (
	cli                    client.Client
	lastNetIOStatTimeStamp int64    //上一次获取net io数据时间点
	lastNetInfo            *NetInfo //上一次net io 数据
	dataBase               string
)

func InitconnInflux(address string, admin string, password string, database string) (err error) {
	dataBase = database
	cli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + address,
		Username: admin,
		Password: password,
	})
	return
}

func getCpuInfo() {
	var cpuInfo = new(CpuInfo)
	// CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
	// fmt.Printf("cpu percent:%v\n", percent)
	// insert
	cpuInfo.CpuPercent = percent[0]
	writesCpuPoints(cpuInfo)
}

// mem info
func getMemInfo() {
	var memInfo = new(MemInfo)
	info, _ := mem.VirtualMemory()
	memInfo.Total = info.Total
	memInfo.Available = info.Available
	memInfo.Used = info.Used
	memInfo.UsedPercent = info.UsedPercent
	memInfo.Buffers = info.Buffers
	memInfo.Cached = info.Cached
	writesMemPoints(memInfo)
}

// disk info
// func getDiskInfo() {
// 	var diskInfo = &DiskInfo{
// 		PartitionUsageStat: make(map[string]*UsageStat, 16),
// 	}
// 	parts, _ := disk.Partitions(true)
// 	for _, part := range parts {
// 		//拿到每一个分区信息
// 		usageStatInfo, err := disk.Usage(part.Mountpoint)
// 		if err != nil {
// 			fmt.Printf("get %s usagestat faild,err:%v", err)
// 			continue
// 		}
// 		usageStat := &UsageStat{
// 			Path:              usageStatInfo.Path,
// 			Fstype:            usageStatInfo.Fstype,
// 			Total:             usageStatInfo.Total,
// 			Used:              usageStatInfo.Used,
// 			UsedPercent:       usageStatInfo.UsedPercent,
// 			InodesTotal:       usageStatInfo.InodesTotal,
// 			InodesUsed:        usageStatInfo.InodesUsed,
// 			InodesFree:        usageStatInfo.InodesFree,
// 			InodesUsedPercent: usageStatInfo.InodesUsedPercent,
// 		}
// 		diskInfo.PartitionUsageStat[part.Mountpoint] = usageStat
// 	}
// 	writesDiskPoints(diskInfo)
// }

func getDiskInfo() {
	var diskInfo = &DiskInfo{
		PartitionUsageStat: make(map[string]*disk.UsageStat, 16),
	}
	parts, _ := disk.Partitions(true)
	for _, part := range parts {
		//拿到每一个分区信息
		usageStat, err := disk.Usage(part.Mountpoint)
		if err != nil {
			fmt.Printf("get %s usagestat faild,err:%v", err)
			continue
		}
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageStat
	}
	writesDiskPoints(diskInfo)
}

func getNetInfo() {
	var netInfo = &NetInfo{
		NetIOCountersStat: make(map[string]*IOStat, 8),
	}
	netIOs, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("get net io counters faild,err:%v", err)
		return
	}
	currentTimeStamp := time.Now().Unix()
	for _, netIO := range netIOs {
		var ioStat = new(IOStat)

		//记录发收包数据
		ioStat.BytesSent = netIO.BytesSent
		ioStat.BytesRecv = netIO.BytesRecv
		ioStat.PacketsSent = netIO.PacketsSent
		ioStat.PacketsRecv = netIO.PacketsRecv

		//将具体网卡数据的ioStat变量添加到Map中
		netInfo.NetIOCountersStat[netIO.Name] = ioStat

		//计算相关速率
		//第一次计算跳过
		if lastNetIOStatTimeStamp == 0 || lastNetInfo == nil {
			continue
		}
		//计算时间间隔
		interval := currentTimeStamp - lastNetIOStatTimeStamp
		//计算速率
		ioStat.BytesSentRate = (float64(ioStat.BytesSent) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].BytesSent)) / float64(interval)
		ioStat.BytesRecvRate = (float64(ioStat.BytesRecv) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].BytesRecv)) / float64(interval)
		ioStat.PacketsSentRate = (float64(ioStat.PacketsSent) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].PacketsSent)) / float64(interval)
		ioStat.PacketsRecvRate = (float64(ioStat.PacketsRecv) - float64(lastNetInfo.NetIOCountersStat[netIO.Name].PacketsRecv)) / float64(interval)

	}
	lastNetIOStatTimeStamp = currentTimeStamp //更新时间
	lastNetInfo = netInfo
	writesNetPoints(netInfo)
}

func Run(interval time.Duration) {
	ticker := time.Tick(interval)
	for _ = range ticker {
		getCpuInfo()
		getMemInfo()
		getDiskInfo()
		getNetInfo()
	}
}
