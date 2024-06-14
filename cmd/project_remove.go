package cmd

import (
	"os"

	"github.com/seabreeze-project/seabreeze/projects"
	"github.com/seabreeze-project/seabreeze/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectRemoveOptions struct {
	Base string
}

var (
	projectRemoveOpt projectRemoveOptions
)

// projectRemoveCmd represents the 'project create' command
var projectRemoveCmd = &cobra.Command{
	Use:     "remove [flags] <name>",
	Aliases: []string{"rm"},
	Short:   "Remove a project",
	Long:    `This command will remove a project.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		h := util.NewHelper(cmd)

		r := projects.NewRepository(viper.GetString("bases.main"))
		project, err := r.Resolve(args[0], projectRemoveOpt.Base)
		if err != nil {
			h.Fatal(err)
		}

		err = os.RemoveAll(project.Path.Get())
		if err != nil {
			h.Fatal(err)
		}

		h.Status("Project removed")
	},
}

func init() {
	projectCmd.AddCommand(projectRemoveCmd)
	rootCmd.AddCommand(projectRemoveCmd)

	projectRemoveCmd.Flags().StringVarP(&projectRemoveOpt.Base, "base", "b", "", "Use given projects base")
}
