package core

import (
	"os"
	"runtime"

	"github.com/docker/docker/client"
)

func createDockerClient() (*client.Client, error) {
	if os.Getenv("DOCKER_HOST") == "" {
		if runtime.GOOS == "windows" {
			os.Setenv("DOCKER_HOST", "npipe:////./pipe/docker_engine")
		} else {
			os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")
		}
	}

	return client.NewClientWithOpts(client.FromEnv)
}
