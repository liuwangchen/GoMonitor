package Operation

import (
	"GoMonitor/Info"

	"github.com/shirou/gopsutil/net"
)

type CpuData []Info.CpuInfo
type NetData []net.IOCountersStat
type ProcessData []Info.ProcessInfo

func (data ProcessData) Clone() ProcessData {
	cloneData := make(ProcessData, len(data))
	copy(cloneData, data)
	return cloneData
}

func (data NetData) Clone() NetData {
	cloneData := make(NetData, len(data))
	copy(cloneData, data)
	return cloneData
}

func (data CpuData) Clone() CpuData {
	cloneData := make(CpuData, len(data))
	copy(cloneData, data)
	return cloneData
}
