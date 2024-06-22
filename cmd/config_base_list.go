package cmd

import (
	"github.com/seabreeze-project/seabreeze/console"
	"github.com/spf13/cobra"
)

// configBaseListCmd represents the 'config base list' command
var configBaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available bases",
	Long:  `This command lists available bases.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		h := console.NewHelper(cmd)

		table := h.Table()
		table.AddRow("NAME", "PATH")
		for name, path := range Core.Config().Bases {
			if name == "main" {
				name = "[main]"
			}
			table.AddRow(name, path)
		}
		h.Println(table)
	},
}

func init() {
	configBaseCmd.AddCommand(configBaseListCmd)
}
