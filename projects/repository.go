package projects

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/seabreeze-project/seabreeze/util"
)

type ProjectRepository struct {
	mainBase *util.Path
}

func NewRepository(mainBase string) *ProjectRepository {
	return &ProjectRepository{
		mainBase: util.NewPath(mainBase),
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

	resolvedBase, err := filepath.Abs(base)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(resolvedBase); os.IsNotExist(err) {
		return nil, fmt.Errorf("projects base directory %q does not exist", base)
	}

	return util.NewPath(resolvedBase), nil
}
