package Operation

import (
	"GoMonitor/Config"
	"GoMonitor/Info"
	"reflect"
	"sort"

	"github.com/shirou/gopsutil/net"
)

type CustomCpuSort struct {
	CpuList  CpuData
	SortFunc func(ci, cj Info.CpuInfo) bool
}

func (cs CustomCpuSort) Len() int {
	return len(cs.CpuList)
}

func (cs CustomCpuSort) Swap(i, j int) {
	cs.CpuList[i], cs.CpuList[j] = cs.CpuList[j], cs.CpuList[i]
}

func (cs CustomCpuSort) Less(i, j int) bool {
	return cs.SortFunc(cs.CpuList[i], cs.CpuList[j])
}

//Sort 根据用户配置来进行Cpu数据排序
func (data CpuData) Sort(conn string) CpuData {
	if sc, ok := Config.OperaCpuConfig[conn]; ok {
		if len(sc.SortPropertyName) > 0 {
			sortData := data.Clone()
			uc := CustomCpuSort{
				CpuList: sortData,
				SortFunc: func(ci, cj Info.CpuInfo) bool {
					v1 := reflect.ValueOf(ci).FieldByName(sc.SortPropertyName)
					if v1.IsValid() {
						v2 := reflect.ValueOf(cj).FieldByName(sc.SortPropertyName)
						return v1.Float() < v2.Float()
					}
					return false
				},
			}
			if sc.Ad == Config.ASC {
				sort.Sort(uc)
			} else {
				sort.Sort(sort.Reverse(uc))
			}
			return sortData
		} else {
			return data
		}
	} else {
		return data
	}
}

type CustomNetSort struct {
	NetList  []net.IOCountersStat
	SortFunc func(ni, nj net.IOCountersStat) bool
}

func (cs CustomNetSort) Len() int {
	return len(cs.NetList)
}

func (cs CustomNetSort) Swap(i, j int) {
	cs.NetList[i], cs.NetList[j] = cs.NetList[j], cs.NetList[i]
}

func (cs CustomNetSort) Less(i, j int) bool {
	return cs.SortFunc(cs.NetList[i], cs.NetList[j])
}

//Sort 根据用户配置来进行Net数据排序
func (data NetData) Sort(conn string) NetData {
	if sc, ok := Config.OperaNetConfig[conn]; ok {
		if len(sc.SortPropertyName) > 0 {
			sortData := data.Clone()
			uc := CustomNetSort{
				NetList: sortData,
				SortFunc: func(ni, nj net.IOCountersStat) bool {
					v1 := reflect.ValueOf(ni).FieldByName(sc.SortPropertyName)
					if v1.IsValid() {
						v2 := reflect.ValueOf(nj).FieldByName(sc.SortPropertyName)
						return v1.Uint() < v2.Uint()
					}
					return false
				},
			}
			if sc.Ad == Config.ASC {
				sort.Sort(uc)
			} else {
				sort.Sort(sort.Reverse(uc))
			}
			return sortData
		} else {
			return data
		}
	} else {
		return data
	}
}

type CustomProcessSort struct {
	ProcessList []Info.ProcessInfo
	SortFunc    func(pi, pj Info.ProcessInfo) bool
}

func (cs CustomProcessSort) Len() int {
	return len(cs.ProcessList)
}

func (cs CustomProcessSort) Swap(i, j int) {
	cs.ProcessList[i], cs.ProcessList[j] = cs.ProcessList[j], cs.ProcessList[i]
}

func (cs CustomProcessSort) Less(i, j int) bool {
	return cs.SortFunc(cs.ProcessList[i], cs.ProcessList[j])
}

//Sort 根据用户配置来进行进程数据排序
func (data ProcessData) Sort(conn string) ProcessData {
	if sc, ok := Config.OperaProcessConfig[conn]; ok {
		if len(sc.SortPropertyName) > 0 {
			sortData := data.Clone()
			uc := CustomProcessSort{
				ProcessList: sortData,
				SortFunc: func(pi, pj Info.ProcessInfo) bool {
					v1 := reflect.ValueOf(pi).FieldByName(sc.SortPropertyName)
					if v1.IsValid() {
						v2 := reflect.ValueOf(pj).FieldByName(sc.SortPropertyName)
						if sc.SortPropertyName == "Id" {
							return v1.Int() < v2.Int()
						}
						return v1.Float() < v2.Float()
					}
					return false
				},
			}
			if sc.Ad == Config.ASC {
				sort.Sort(uc)
			} else {
				sort.Sort(sort.Reverse(uc))
			}
			return sortData
		} else {
			return data
		}
	} else {
		return data
	}
}
