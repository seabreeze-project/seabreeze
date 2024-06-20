package cmd

import (
	"path/filepath"

	"github.com/seabreeze-project/seabreeze/projects"
	"github.com/seabreeze-project/seabreeze/util"
	"github.com/spf13/cobra"
)

type projectInitOptions struct {
	projectCreateOptions
}

var (
	projectInitOpt projectInitOptions
)

// projectInitCmd represents the 'project init' command
var projectInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create new project in current directory",
	Long:  `This command will create a new project in the current directory.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		h := util.NewHelper(cmd)

		dir, err := Core.Here()
		if err != nil {
			h.Fatal(err)
		}

		h.Printf("Initializing new project in %q...\n", dir)

		createOpt := projects.CreateOptions{
			ProjectManifest: &projects.ProjectMetadata{
				Description: projectCreateOpt.Description,
			},
			AppConfig: Core.Config(),
		}

		templateFile, err := filepath.Abs(projectCreateOpt.TemplateFile)
		if err != nil {
			h.Fatal(err)
		}
		createOpt.TemplateFile = templateFile
		h.Println("  Using template:", templateFile)

		_, err = projects.Create(dir, createOpt)
		if err != nil {
			h.Fatal(err)
		}

		h.Status("Project initialized successfully.")
	},
}

func init() {
	projectCmd.AddCommand(projectInitCmd)
	rootCmd.AddCommand(projectInitCmd)

	projectInitCmd.Flags().StringVarP(&projectInitOpt.Description, "description", "d", "", "Set project description")
	projectInitCmd.Flags().StringVarP(&projectInitOpt.TemplateFile, "template", "t", Core.ConfigPath("docker-compose.yml.tmpl"), "Use given template file to create the project compose file")
}
