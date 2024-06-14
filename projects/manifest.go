package projects

import (
	"errors"
	"fmt"
	"os"

	"github.com/seabreeze-project/seabreeze/util"
)

var ManifestFileName = "seabreeze.%s"

var ErrManifestNotFound = errors.New("project manifest not found")

func FindManifest(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", errors.New("provided path is not a directory")
	}

	p := util.NewPath(path)
	filePath, exists := p.Locate(fmt.Sprintf(ManifestFileName, "yml"), fmt.Sprintf(ManifestFileName, "yaml"))
	if !exists {
		return "", ErrManifestNotFound
	}

	return filePath, nil
}

func ManifestExists(path string) (bool, error) {
	_, err := FindManifest(path)
	if errors.Is(err, ErrManifestNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
