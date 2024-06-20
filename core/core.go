package core

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

type Core struct {
	ConfigBasePath string
	viper          *viper.Viper
	client         *client.Client
	config         *Configuration
}

func New() *Core {
	configBasePath := "/etc/seabreeze"
	if runtime.GOOS == "windows" {
		configBasePath = os.Getenv("PROGRAMDATA") + "\\Seabreeze"
	} else if runtime.GOOS == "darwin" {
		configBasePath = "/Library/Application Support/Seabreeze"
	}

	return &Core{
		ConfigBasePath: configBasePath,
		viper: viper.New(),
	}
}

func (c *Core) Here() (string, error) {
	return os.Getwd()
}

func (c *Core) ConfigPath(path string) string {
	return filepath.Join(c.ConfigBasePath, path)
}

func (c *Core) Config() *Configuration {
	if c.config == nil {
		panic("accessing config before it has been loaded")
	}
	return c.config
}

func (c *Core) ConfigFileUsed() string {
	return c.viper.ConfigFileUsed()
}

func (c *Core) LoadConfig(path string) error {
	if path != "" {
		c.viper.SetConfigFile(path)
	} else {
		c.viper.AddConfigPath(c.ConfigBasePath)
		c.viper.SetConfigType("yml")
		c.viper.SetConfigName("config")
	}

	c.viper.SetEnvPrefix("SEABREEZE")
	c.viper.AutomaticEnv()

	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	c.config = &Configuration{}
	if err := c.viper.Unmarshal(&c.config); err != nil {
		return err
	}

	return nil
}

func (c *Core) Client() (*client.Client, error) {
	if c.client != nil {
		return c.client, nil
	}

	if os.Getenv("DOCKER_HOST") == "" {
		if runtime.GOOS == "windows" {
			os.Setenv("DOCKER_HOST", "npipe:////./pipe/docker_engine")
		} else {
			os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	c.client = cli
	return cli, nil
}
