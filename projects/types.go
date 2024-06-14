package projects

type Container struct {
	Name string
}

type ProjectStatus struct {
	Online int
	Total  int
}

type CreateOptions struct {
	ProjectManifest *ProjectMetadata
	TemplateFile    string
}
