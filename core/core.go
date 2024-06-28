package core

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/seabreeze-project/seabreeze/config"
	"github.com/spf13/viper"
)

// Core gives access to essential parts of the application.
type Core struct {
	ConfigBasePath string
	client         *client.Client
	config         *config.Configuration
}

// New returns a new Core instance.
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

// Here returns the current working directory.
func (c *Core) Here() (string, error) {
	return os.Getwd()
}

// ConfigPath returns the path to the given config file.
func (c *Core) ConfigPath(path string) string {
	return filepath.Join(c.ConfigBasePath, path)
}

// Config returns the loaded configuration.
func (c *Core) Config() *config.Configuration {
	if c.config == nil {
		panic("accessing config before it has been loaded")
	}
	return c.config
}

// LoadConfig loads the configuration from the given path or from the default location if empty.
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

	config, err := config.LoadConfig(v)
	if err != nil {
		return err
	}

	c.config = config
	return nil
}

// Client returns a preconfigured Docker API client.
func (c *Core) Client() (*client.Client, error) {
	if c.client != nil {
		return c.client, nil
	}

	cli, err := createDockerClient()
	if err != nil {
		return nil, err
	}

	c.client = cli
	return cli, nil
}
