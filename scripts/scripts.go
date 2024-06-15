package scripts

import (
	"fmt"
	"io"
	"os/exec"
	"os/user"

	"github.com/seabreeze-project/seabreeze/projects"
)

type Script struct {
	Name    string
	Config  *projects.ScriptConfig
	Project *projects.Project
}

func New(name string, config *projects.ScriptConfig) *Script {
	return &Script{
		Name:   name,
		Config: config,
	}
}

func FromProject(p *projects.Project, scriptName string) (*Script, error) {
	sc, ok := p.Metadata.Scripts[scriptName]
	if !ok {
		return nil, fmt.Errorf("script %q not found", scriptName)
	}

	s := New(scriptName, sc)
	s.Project = p
	return s, nil
}

func (s *Script) Run(opt ScriptRunOptions, args ...string) (*exec.Cmd, error) {
	var (
		cmdline []string
		process *exec.Cmd
		u       *user.User
	)
	sc := s.Config

	target, targetName, err := parseTarget(sc.Target)
	if err != nil {
		return nil, err
	}

	envMap := make(map[string]string)
	if sc.Workdir != "" {
		envMap["PWD"] = sc.Workdir
	}
	if sc.User != "" {
		u, err = user.Lookup(sc.User)
		if err != nil {
			return nil, err
		}
		envMap["HOME"] = u.HomeDir
	}
	if sc.Environment != nil {
		for k, v := range sc.Environment {
			envMap[k] = v
		}
	}

	realCmdline, err := generateCmdline(sc.Command, envMap, args...)
	if err != nil {
		return nil, err
	}

	switch target {
	case "service", "container":
		// TODO: use Docker API instead
		if target == "service" {
			cmdline = append(cmdline, "docker", "compose", "exec")
		} else {
			cmdline = append(cmdline, "docker", "exec")
		}
		if sc.Workdir != "" {
			cmdline = append(cmdline, "-w", sc.Workdir)
		}
		if sc.User != "" {
			cmdline = append(cmdline, "-u", sc.User)
		}
		cmdline = append(cmdline, targetName)
		cmdline = append(cmdline, realCmdline...)
	case "host":
		cmdline = append(cmdline, realCmdline...)
	default:
		return nil, fmt.Errorf("unknown target %q", target)
	}

	process = exec.Command(cmdline[0])
	if len(cmdline) > 1 {
		process.Args = append(process.Args, cmdline[1:]...)
	}

	if target == "host" {
		if sc.Workdir != "" {
			process.Dir = sc.Workdir
		}
		if u != nil {
			err = setupHostProcessForUser(process, u)
			if err != nil {
				return nil, err
			}
		}
	}

	for k, v := range envMap {
		process.Env = append(process.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if opt.Stdout != nil {
		process.Stdout = opt.Stdout
	}
	if opt.Stderr != nil {
		process.Stderr = opt.Stderr
	}

	err = process.Run()
	if err != nil {
		return nil, err
	}

	return process, nil
}

type ScriptRunOptions struct {
	Stdout io.Writer
	Stderr io.Writer
}
