package projects

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/seabreeze-project/seabreeze/util"
	"gopkg.in/yaml.v3"
)

type Project struct {
	ID       string
	Name     string
	Path     *util.Path
	Metadata *ProjectMetadata
}

func Open(path string) (*Project, error) {
	manifest, err := FindManifest(path)
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(manifest)
	if err != nil {
		return nil, err
	}

	metadata := &ProjectMetadata{}
	err = yaml.Unmarshal(file, metadata)
	if err != nil {
		return nil, err
	}

	p := &Project{
		ID:       util.Sha1(path),
		Name:     filepath.Base(path),
		Path:     util.NewPath(path),
		Metadata: metadata,
	}
	return p, nil
}

func OpenHere() (*Project, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return Open(path)
}

func Create(path string, opt CreateOptions) (*Project, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	exists, err := ManifestExists(path)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("project config file already exists")
	}

	var metadata *ProjectMetadata
	if opt.ProjectManifest != nil {
		metadata = opt.ProjectManifest
	} else {
		metadata = &ProjectMetadata{}
	}

	file, err := yaml.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	if err = os.WriteFile(filepath.Join(path, "seabreeze.yml"), file, 0644); err != nil {
		return nil, err
	}

	p := &Project{
		ID:       util.Sha1(path),
		Name:     filepath.Base(path),
		Path:     util.NewPath(path),
		Metadata: metadata,
	}

	if err := p.prepare(opt); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Project) ComposeConfig() (*ComposeConfig, error) {
	file, exists := p.Path.Locate("docker-compose.yml", "docker-compose.yaml")
	if !exists {
		return nil, errors.New("compose file not found")
	}

	return ParseDockerCompose(file)
}

func (p *Project) Container(service string, replica int) (*Container, error) {
	c, err := p.ComposeConfig()
	if err != nil {
		return nil, err
	}

	s, ok := c.Services[service]
	if !ok {
		return nil, fmt.Errorf("service %q not found", service)
	}

	var containerName string
	if s.ContainerName != "" {
		if replica > 1 {
			return nil, fmt.Errorf("cannot access replica %d of service %q: service does not support replicas because it has a static container name", replica, service)
		}
		containerName = s.ContainerName
	} else {
		if replica > s.Replicas {
			return nil, fmt.Errorf("cannot access replica %d of service %q: service only defines %d replicas", replica, service, s.Replicas)
		}
		containerName = fmt.Sprintf("%s-%s-%d", p.Name, service, replica)
	}

	return &Container{
		Name: containerName,
	}, nil
}

func (p *Project) Status(cli *client.Client) (*ProjectStatus, error) {
	cp, err := p.ComposeConfig()
	if err != nil {
		return nil, err
	}

	var online int
	var total int

	// TODO: This is just quick and dirty, and thus an unreliable solution. Use an approach that relies on container labels.
	for serviceName := range cp.Services {
		ct, err := p.Container(serviceName, 1)
		if err != nil {
			return nil, err
		}

		container, err := cli.ContainerInspect(context.Background(), ct.Name)
		if err != nil {
			return nil, err
		}

		if container.State.Running {
			online++
		}
		total++
	}

	return &ProjectStatus{
		Online: online,
		Total:  total,
	}, nil
}

func (p *Project) prepare(opt CreateOptions) error {
	templateFile := opt.TemplateFile

	templateVariables := map[string]any{
		"Project": p,
	}
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}
	f, err := os.Create(p.Path.Join("docker-compose.yml"))
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, templateVariables)
	if err != nil {
		return err
	}

	return nil
}
