package cmd

import (
	"github.com/spf13/cobra"
)

// projectCmd represents the 'project' command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  `Manage projects`,
	Args:  projectListCmd.Args,
	Run:   projectListCmd.Run,
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
