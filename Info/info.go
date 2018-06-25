package Info

import (
	"time"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type CpuInfo struct {
	Name string  `json:"name,omitempty"`
	Used float64 `json:"used,omitempty"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total,omitempty"`
	Available   uint64  `json:"available,omitempty"`
	Used        uint64  `json:"used,omitempty"`
	UsedPercent float64 `json:"used_percent,omitempty"`
}

type ProcessInfo struct {
	MemoryPercent float32 `json:"memory_percent"`
	Name          string  `json:"name"`
	Id            int32   `json:"id"`
	CPUPercent    float64 `json:"cpu_percent"`
	Status        string  `json:"status"`
}

func GetCpuInfo() []CpuInfo {
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	cpuList := make([]CpuInfo, 0, len(c))
	for i, _ := range c {
		cpuList = append(cpuList, CpuInfo{
			Name: c[i].ModelName,
			Used: cc[i],
		})
	}
	return cpuList
}

func GetMemoryInfo() MemoryInfo {
	v, _ := mem.VirtualMemory()
	return MemoryInfo{
		Total:       v.Total / 1024 / 1024,
		Available:   v.Available / 1024 / 1024,
		Used:        v.Used / 1024 / 1024,
		UsedPercent: v.UsedPercent,
	}
}

func GetNetInfo() []net.IOCountersStat {
	nv, _ := net.IOCounters(true)
	return nv	
}

func GetProcessInfo() []ProcessInfo {
	p, _ := process.Processes()
	processList := make([]ProcessInfo, 0)
	for _, pchild := range p {
		if pchild.Pid != 0 {
			memoryPercent, _ := pchild.MemoryPercent()
			name, _ := pchild.Name()
			status, err := pchild.Status()
			if err != nil {
				status = "不支持"
			}
			cpuPercent, _ := pchild.CPUPercent()
			processList = append(processList, ProcessInfo{
				MemoryPercent: memoryPercent,
				Name:          name,
				Id:            pchild.Pid,
				Status:        status,
				CPUPercent:    cpuPercent,
			})
		}
	}
	return processList
}
