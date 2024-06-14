package cmd

import (
	"fmt"

	"github.com/seabreeze-project/seabreeze/projects"
	"github.com/seabreeze-project/seabreeze/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectListOptions struct {
	Base string
}

var (
	projectListOpt projectListOptions
)

// projectListCmd represents the 'project list' command
var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long:  `List available projects`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		h := util.NewHelper(cmd)

		r := projects.NewRepository(viper.GetString("bases.main"))
		list, err := r.List(projectListOpt.Base)
		if err != nil {
			h.Fatal(err)
		}

		table := h.Table()
		table.AddRow("ID", "NAME", "DESCRIPTION", "DIRECTORY", "STATUS")

		for _, projectPath := range list {
			cli, err := Core.Client()
			if err != nil {
				h.Fatal(err)
			}

			p, err := projects.Open(projectPath.Get())
			if err != nil {
				h.Fatal(err)
			}
			status, err := p.Status(cli)
			if err != nil {
				h.Fatal(err)
			}

			table.AddRow(
				p.ID[0:12],
				p.Name,
				p.Metadata.Description,
				projectPath.Dir(),
				fmt.Sprintf("%d online / %d", status.Online, status.Total),
			)
		}

		h.Println(table)
	},
}

func init() {
	projectCmd.AddCommand(projectListCmd)
	rootCmd.AddCommand(projectListCmd)

	projectListCmd.Flags().StringVarP(&projectListOpt.Base, "base", "b", "", "Use given projects base")
}
