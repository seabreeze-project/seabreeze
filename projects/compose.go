package projects

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ComposeConfig struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Hostname      string   `yaml:"hostname"`
	ContainerName string   `yaml:"container_name"`
	Image         string   `yaml:"image"`
	Ports         []string `yaml:"ports"`
	Replicas      int      `yaml:"replicas"`
}

func ParseDockerCompose(path string) (*ComposeConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &ComposeConfig{}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
