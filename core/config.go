package core

import (
	"errors"
)

type Configuration struct {
	Bases        BasesMap           `mapstructure:"bases"`
	ReverseProxy ReverseProxyConfig `mapstructure:"reverse_proxy"`
}

type BasesMap map[string]string

type ReverseProxyConfig struct {
	NetworkName string `mapstructure:"network_name"`
}

func (c *Configuration) Validate() error {
	if _, ok := c.Bases["main"]; !ok {
		return errors.New("main base is not defined")
	}

	return nil
}
