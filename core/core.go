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

func (c *Core) LoadConfig(path string) error {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.AddConfigPath(c.ConfigBasePath)
		v.SetConfigType("yml")
		v.SetConfigName("config")
	}

	v.SetEnvPrefix("SEABREEZE")
	v.AutomaticEnv()

	config, err := LoadConfig(v)
	if err != nil {
		return err
	}

	c.config = config
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
