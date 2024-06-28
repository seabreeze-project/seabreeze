package core

import (
	"os"
	"runtime"

	"github.com/docker/docker/client"
)

func createDockerClient() (*client.Client, error) {
	// Prevent client and API version mismatch
	os.Setenv("DOCKER_API_VERSION", "1.45")

	if os.Getenv("DOCKER_HOST") == "" {
		if runtime.GOOS == "windows" {
			os.Setenv("DOCKER_HOST", "npipe:////./pipe/docker_engine")
		} else {
			os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")
		}
	}

	return client.NewClientWithOpts(client.FromEnv)
}
