package core

type Configuration struct {
	Bases        BasesMap           `mapstructure:"bases"`
	ReverseProxy ReverseProxyConfig `mapstructure:"reverse_proxy"`
}

type BasesMap map[string]string

type ReverseProxyConfig struct {
	NetworkName string `mapstructure:"network_name"`
}
