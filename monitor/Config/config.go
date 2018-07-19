package Config

import (
	"strings"
)

const (
	ASC = iota
	DESC
)

type OperationConfig struct {
	SearchStr        string
	SortPropertyName string
	Ad               int
}

type operaConfigMap map[string]*OperationConfig

//排序配置
var OperaCpuConfig operaConfigMap = make(map[string]*OperationConfig)

//排序配置
var OperaNetConfig operaConfigMap = make(map[string]*OperationConfig)

//排序配置
var OperaProcessConfig operaConfigMap = make(map[string]*OperationConfig)

func (scm operaConfigMap) setSortConfig(conn, propertyName string) {
	if len(strings.TrimSpace(conn)) > 0 {
		if config, ok := scm[conn]; ok {
			if config.SortPropertyName != propertyName {
				config.Ad = DESC
			} else {
				if config.Ad == ASC {
					config.Ad = DESC
				} else {
					config.Ad = ASC
				}
			}
			config.SortPropertyName = propertyName
		} else {
			scm[conn] = &OperationConfig{
				SortPropertyName: propertyName,
				Ad:               DESC,
			}
		}
	}
}

func SetCpuSortConfig(conn, propertyName string) {
	OperaCpuConfig.setSortConfig(conn, propertyName)
}

func SetNetSortConfig(conn, propertyName string) {
	OperaNetConfig.setSortConfig(conn, propertyName)
}

func SetProcessSortConfig(conn, propertyName string) {
	OperaProcessConfig.setSortConfig(conn, propertyName)
}

func (scm operaConfigMap) setSearchConfig(conn, searchStr string) {
	if len(strings.TrimSpace(conn)) > 0 {
		if config, ok := scm[conn]; ok {
			config.SearchStr = searchStr
		} else {
			scm[conn] = &OperationConfig{
				SearchStr: searchStr,
			}
		}
	}
}

func SetCpuSearchConfig(conn, searchStr string) {
	OperaCpuConfig.setSearchConfig(conn, searchStr)
}

func SetNetSearchConfig(conn, searchStr string) {
	OperaNetConfig.setSearchConfig(conn, searchStr)
}

func SetProcessSearchConfig(conn, searchStr string) {
	OperaProcessConfig.setSearchConfig(conn, searchStr)
}
