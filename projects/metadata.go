package projects

type ProjectMetadata struct {
	Description string                   `yaml:"description"`
	Scripts     map[string]*ScriptConfig `yaml:"scripts"`
}

type ScriptConfig struct {
	Command     string            `yaml:"command"`
	Workdir     string            `yaml:"workdir"`
	Container   string            `yaml:"container"`
	User        string            `yaml:"user"`
	Target      string            `yaml:"target"`
	Environment map[string]string `yaml:"environment"`
}
