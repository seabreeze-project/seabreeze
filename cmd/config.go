package cmd

import (
	"github.com/seabreeze-project/seabreeze/console"
	"github.com/spf13/cobra"
)

// configCmd represents the 'config' command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the configuration",
	Long: `This command group allows you to manage the configuration.
Calling the "config" command without a subcommand will print information about the loaded configuration.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		h := console.NewHelper(cmd)

		table := h.Table()
		table.Wrap = true
		table.AddRow("Config file used:", Core.Config().SourceFile())
		h.Println(table)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
