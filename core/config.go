package core

type Configuration struct {
	Bases        BasesConfig        `mapstructure:"bases"`
	ReverseProxy ReverseProxyConfig `mapstructure:"reverse_proxy"`
}

type BasesConfig struct {
	Main string `mapstructure:"main"`
}

type ReverseProxyConfig struct {
	NetworkName string `mapstructure:"network_name"`
}
