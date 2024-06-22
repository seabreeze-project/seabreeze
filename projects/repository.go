package projects

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/seabreeze-project/seabreeze/core"
	"github.com/seabreeze-project/seabreeze/util"
)

type ProjectRepository struct {
	bases    core.BasesMap
	mainBase *util.Path
}

func NewRepository(bases core.BasesMap) *ProjectRepository {
	if bases == nil {
		panic("bases map must not be nil")
	}

	mainBasePath, ok := bases["main"]
	if !ok {
		panic("main base is not defined")
	}

	return &ProjectRepository{
		bases:    bases,
		mainBase: util.NewPath(mainBasePath),
	}
}

func (r *ProjectRepository) List(base string) ([]*util.Path, error) {
	var projects []*util.Path

	basePath, err := r.ResolveBase(base)
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(basePath.Get(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		exists, err := ManifestExists(path)
		if err != nil {
			return err
		}
		if !exists {
			return nil
		}

		projects = append(projects, util.NewPath(path))
		return nil
	})

	return projects, err
}

func (r *ProjectRepository) Exists(name string, base string) (bool, error) {
	basePath, err := r.ResolveBase(base)
	if err != nil {
		return false, err
	}

	return ManifestExists(basePath.Join(name))
}

func (r *ProjectRepository) Open(name string) (*Project, error) {
	return Open(r.mainBase.Join(name))
}

func (r *ProjectRepository) Resolve(name string, base string) (*Project, error) {
	if name == "." && base == "" {
		return OpenHere()
	}

	basePath, err := r.ResolveBase(base)
	if err != nil {
		return nil, err
	}

	resolvedProject := basePath.Join(name)
	if _, err := os.Stat(resolvedProject); os.IsNotExist(err) {
		return nil, fmt.Errorf("project %q does not exist", name)
	}

	return Open(resolvedProject)
}

func (r *ProjectRepository) ResolveBase(base string) (*util.Path, error) {
	if base == "" {
		return r.mainBase, nil
	}

	var resolvedBase string
	if base[0] == '@' {
		baseName := base[1:]
		if baseName == "" {
			return nil, fmt.Errorf("base name cannot be empty")
		}

		var ok bool
		resolvedBase, ok = r.bases[baseName]
		if !ok {
			return nil, fmt.Errorf("unknown base name %q", baseName)
		}
	} else {
		var err error
		resolvedBase, err = filepath.Abs(base)
		if err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(resolvedBase); os.IsNotExist(err) {
		return nil, fmt.Errorf("projects base directory %q does not exist", base)
	}

	return util.NewPath(resolvedBase), nil
}
