package core

import (
	"errors"

	"github.com/spf13/viper"
)

type Configuration struct {
	provider     *viper.Viper       `mapstructure:"-"`
	Bases        BasesMap           `mapstructure:"bases"`
	ReverseProxy ReverseProxyConfig `mapstructure:"reverse_proxy"`
}

type BasesMap map[string]string

type ReverseProxyConfig struct {
	NetworkName string `mapstructure:"network_name"`
}

func (c *Configuration) Get(key string) any {
	return c.provider.Get(key)
}

func (c *Configuration) Validate() error {
	if _, ok := c.Bases["main"]; !ok {
		return errors.New("main base is not defined")
	}

	return nil
}
