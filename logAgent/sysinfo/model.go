package sysinfo

import (
	"github.com/shirou/gopsutil/disk"
)

type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}
type MemInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedpercent"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
}

type UsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"cached"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"used_Percent"`
	InodesTotal       uint64  `json:"inodes_total"`
	InodesUsed        uint64  `json:"inodes_used"`
	InodesFree        uint64  `json:"inodes_free"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
}

type DiskInfo struct {
	PartitionUsageStat map[string]*disk.UsageStat
}

type IOStat struct {
	BytesSent       uint64
	BytesRecv       uint64
	PacketsSent     uint64
	PacketsRecv     uint64
	BytesSentRate   float64 `json:"bytes_sent_rate"`
	BytesRecvRate   float64 `json:"bytes_recv_rate"`
	PacketsSentRate float64 `json:"packets_sent_rate"`
	PacketsRecvRate float64 `json:"packets_recv_rate"`
}

type NetInfo struct {
	NetIOCountersStat map[string]*IOStat
}
