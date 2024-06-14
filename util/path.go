package util

import (
	"os"
	"path/filepath"
)

type Path struct {
	base string
}

func NewPath(path string) *Path {
	return &Path{path}
}

func (p *Path) Get() string {
	return p.base
}

func (p *Path) String() string {
	return p.base
}

func (p *Path) Sub(name string) *Path {
	return NewPath(p.Join(name))
}

func (p *Path) Join(parts ...string) string {
	args := append([]string{p.base}, parts...)
	return filepath.Join(args...)
}

func (p *Path) Locate(names ...string) (string, bool) {
	for _, name := range names {
		path := p.Join(name)
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}
	return "", false
}

func (p *Path) Dir() *Path {
	return NewPath(filepath.Dir(p.base))
}
