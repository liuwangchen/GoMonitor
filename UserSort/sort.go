package UserSort

import (
	"GoMonitor/Info"
	"reflect"
	"sort"

	"github.com/shirou/gopsutil/net"
)

type CpuData []Info.CpuInfo
type NetData []net.IOCountersStat
type ProcessData []Info.ProcessInfo

type SortConfig struct {
	PropertyName string
	Ad           string
}

//排序配置
type sortConfigMap map[string]*SortConfig

//排序配置
var SortCpuConfig sortConfigMap = make(map[string]*SortConfig)

//排序配置
var SortNetConfig sortConfigMap = make(map[string]*SortConfig)

//排序配置
var SortProcessConfig sortConfigMap = make(map[string]*SortConfig)

func (scm sortConfigMap) setSortConfig(address, propertyName, ad string) {
	if config, ok := scm[address]; ok {
		config.PropertyName = propertyName
		config.Ad = ad
	} else {
		scm[address] = &SortConfig{
			PropertyName: propertyName,
			Ad:           ad,
		}
	}
}

func SetCpuSortConfig(address, propertyName, ad string) {
	SortCpuConfig.setSortConfig(address, propertyName, ad)
}

func SetNetSortConfig(address, propertyName, ad string) {
	SortNetConfig.setSortConfig(address, propertyName, ad)
}

func SetProcessSortConfig(address, propertyName, ad string) {
	SortProcessConfig.setSortConfig(address, propertyName, ad)
}

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
func (data CpuData) Sort(addrs string) CpuData {
	if sc, ok := SortCpuConfig[addrs]; ok {
		uc := CustomCpuSort{
			CpuList: data,
			SortFunc: func(ci, cj Info.CpuInfo) bool {
				v1 := reflect.ValueOf(ci).FieldByName(sc.PropertyName)
				if v1.IsValid() {
					v2 := reflect.ValueOf(cj).FieldByName(sc.PropertyName)
					return v1.Float() < v2.Float()
				}
				return false
			},
		}
		if sc.Ad == "asc" {
			sort.Sort(uc)
		} else {
			sort.Sort(sort.Reverse(uc))
		}
	}
	return data
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
func (data NetData) Sort(addrs string) NetData {
	if sc, ok := SortNetConfig[addrs]; ok {
		uc := CustomNetSort{
			NetList: data,
			SortFunc: func(ni, nj net.IOCountersStat) bool {
				v1 := reflect.ValueOf(ni).FieldByName(sc.PropertyName)
				if v1.IsValid() {
					v2 := reflect.ValueOf(nj).FieldByName(sc.PropertyName)
					return v1.Uint() < v2.Uint()
				}
				return false
			},
		}
		if sc.Ad == "asc" {
			sort.Sort(uc)
		} else {
			sort.Sort(sort.Reverse(uc))
		}
	}
	return data
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
func (data ProcessData) Sort(addrs string) ProcessData {
	if sc, ok := SortProcessConfig[addrs]; ok {
		uc := CustomProcessSort{
			ProcessList: data,
			SortFunc: func(pi, pj Info.ProcessInfo) bool {
				v1 := reflect.ValueOf(pi).FieldByName(sc.PropertyName)
				if v1.IsValid() {
					v2 := reflect.ValueOf(pj).FieldByName(sc.PropertyName)
					if sc.PropertyName == "Id" {
						return v1.Int() < v2.Int()
					}
					return v1.Float() < v2.Float()
				}
				return false
			},
		}
		if sc.Ad == "asc" {
			sort.Sort(uc)
		} else {
			sort.Sort(sort.Reverse(uc))
		}
	}
	return data
}
