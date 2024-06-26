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

func LoadConfig(v *viper.Viper) (*Configuration, error) {
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Configuration{provider: v}
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Configuration) Get(key string) any {
	return c.provider.Get(key)
}

func (c *Configuration) SourceFile() string {
	return c.provider.ConfigFileUsed()
}

func (c *Configuration) Validate() error {
	if _, ok := c.Bases["main"]; !ok {
		return errors.New("main base is not defined")
	}

	return nil
}
