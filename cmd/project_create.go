package cmd

import (
	"path/filepath"

	"github.com/seabreeze-project/seabreeze/projects"
	"github.com/seabreeze-project/seabreeze/util"
	"github.com/spf13/cobra"
)

type projectCreateOptions struct {
	Base         string
	Description  string
	TemplateFile string
}

var (
	projectCreateOpt projectCreateOptions
)

// projectCreateCmd represents the 'project create' command
var projectCreateCmd = &cobra.Command{
	Use:   "create [flags] <name>",
	Short: "Create new project",
	Long:  `This command will create a new project.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		h := util.NewHelper(cmd)

		r := projects.NewRepository(Core.Config().Bases.Main)
		base, err := r.ResolveBase(projectCreateOpt.Base)
		if err != nil {
			h.Fatal("cannot open projects base:", err)
		}

		projectName := args[0]
		projectDir := base.Join(projectName)

		h.Printf("Creating new project %q...\n", projectName)

		createOpt := projects.CreateOptions{
			ProjectManifest: &projects.ProjectMetadata{
				Description: projectCreateOpt.Description,
			},
		}

		templateFile, err := filepath.Abs(projectCreateOpt.TemplateFile)
		if err != nil {
			h.Fatal(err)
		}
		createOpt.TemplateFile = templateFile
		h.Println("  Using template:", templateFile)

		_, err = projects.Create(projectDir, createOpt)
		if err != nil {
			h.Fatal(err)
		}

		h.Status("Project created successfully.")
	},
}

func init() {
	projectCmd.AddCommand(projectCreateCmd)
	rootCmd.AddCommand(projectCreateCmd)

	projectCreateCmd.Flags().StringVarP(&projectCreateOpt.Base, "base", "b", "", "Use given projects base")
	projectCreateCmd.Flags().StringVarP(&projectCreateOpt.Description, "description", "d", "", "Set project description")
	projectCreateCmd.Flags().StringVarP(&projectCreateOpt.TemplateFile, "template", "t", Core.ConfigPath("docker-compose.yml.tmpl"), "Use given template file to create the project compose file")
}
