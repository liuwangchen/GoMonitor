package Config

type cpuConfig struct {
	Sort string
}

type netConfig struct {
}

type processConfig struct {
}

func GetCpuConfig() *cpuConfig {
	return &cpuConfig{}
}

func GetNetConfig() *netConfig {
	return &netConfig{}
}

func GetProcessConfig() *processConfig {
	return &processConfig{}
}
