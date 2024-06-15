package cmd

import (
	"os"

	"github.com/seabreeze-project/seabreeze/projects"
	"github.com/seabreeze-project/seabreeze/scripts"
	"github.com/seabreeze-project/seabreeze/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type scriptOptions struct {
	ProjectsBase string
}

var (
	scriptOpt scriptOptions
)

// scriptCmd represents the 'script' command
var scriptCmd = &cobra.Command{
	Use:   "script [flags] <project-name> <script-name> [-- args...]",
	Args:  cobra.MinimumNArgs(2),
	Short: "Run a defined script",
	Long:  `Run a script defined in the project's manifest file.`,
	Run: func(cmd *cobra.Command, args []string) {
		h := util.NewHelper(cmd)

		projectName := args[0]
		scriptName := args[1]
		scriptArgs := args[2:]

		r := projects.NewRepository(viper.GetString("projects_dir"))
		project, err := r.Resolve(projectName, scriptOpt.ProjectsBase)
		if err != nil {
			h.Fatal(err)
		}

		s, err := scripts.FromProject(project, scriptName)
		if err != nil {
			h.Fatal(err)
		}
		opt := scripts.ScriptRunOptions{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
		proc, err := s.Run(opt, scriptArgs...)
		if err != nil {
			h.Fatal(err)
		}

		os.Exit(proc.ProcessState.ExitCode())
	},
}

func init() {
	rootCmd.AddCommand(scriptCmd)

	scriptCmd.Flags().StringVarP(&scriptOpt.ProjectsBase, "base", "b", "", "Use given projects base")
	scriptCmd.PersistentFlags().StringP("shell", "s", "bash", "The shell to use for the script")
}
