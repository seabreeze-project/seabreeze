package cmd

import (
	"github.com/seabreeze-project/seabreeze/projects"
	"github.com/seabreeze-project/seabreeze/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectEditOptions struct {
	Base string
}

var (
	projectEditOpt projectEditOptions
)

// projectEditCommand represents the 'project edit' command
var projectEditCmd = &cobra.Command{
	Use:   "edit [flags] <name>",
	Short: "Edit a project",
	Long:  `This command will edit a project.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		h := util.NewHelper(cmd)

		r := projects.NewRepository(viper.GetString("bases.main"))
		project, err := r.Resolve(args[0], projectEditOpt.Base)
		if err != nil {
			h.Fatal(err)
		}

		file, exists := project.Path.Locate("docker-compose.yml", "docker-compose.yaml")
		if !exists {
			h.Fatal("No compose file found")
		}

		util.OpenEditor(file)
	},
}

func init() {
	projectCmd.AddCommand(projectEditCmd)
	rootCmd.AddCommand(projectEditCmd)

	projectEditCmd.Flags().StringVarP(&projectEditOpt.Base, "base", "b", "", "Use given projects base")
}
