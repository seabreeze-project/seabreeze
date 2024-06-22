package cmd

import (
	"github.com/spf13/cobra"
)

// configBaseCmd represents the 'config base' command
var configBaseCmd = &cobra.Command{
	Use:   "base",
	Short: "Manage project bases",
	Long: `This command group allows you to manage project bases.
Calling the "base" command without a subcommand will list available bases.`,
	Args: cobra.NoArgs,
	Run: configBaseListCmd.Run,
}

func init() {
	configCmd.AddCommand(configBaseCmd)
	rootCmd.AddCommand(configBaseCmd)
}
