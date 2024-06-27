package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/seabreeze-project/seabreeze/shell"
	"github.com/spf13/cobra"
)

// shellCmd represents the 'shell' command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start Seabreeze in Shell mode",
	Long:  `Starts an interactive Seabreeze shell session.`,
	Run: func(cmd *cobra.Command, args []string) {
		sh := shell.New("Seabreeze Shell").SetAppVersion("0.0.1")

		sh.SetPromptPrefix("seabreeze> ")
		sh.AddCommand(
			&shell.Command{
				Name: "bash",
				Desc: "Run bash",
				Run: func(args []string) error {
					cmd := exec.Command("bash")
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					return cmd.Run()
				},
			},
			&shell.Command{
				Name: "help",
				Desc: "Display help",
				Run: func(args []string) error {
					fmt.Println(sh.Help())
					return nil
				},
			},
			&shell.Command{
				Name:    "clear",
				Desc:    "Clear the screen",
				Aliases: []string{"cls"},
				Run: func(args []string) error {
					cmd := exec.Command("clear")
					if runtime.GOOS == "windows" {
						cmd = exec.Command("cmd", "/c", "cls")
					}
					cmd.Stdout = os.Stdout
					return cmd.Run()
				},
			},
			&shell.Command{
				Name:    "exit",
				Desc:    "Exit the shell",
				Aliases: []string{"quit"},
				Run: func(args []string) error {
					os.Exit(0)
					return nil
				},
			},
		)
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == "shell" || cmd.Name() == "help" {
				continue
			}
			sh.AddExternalCommand(cmd) // TODO: Find an alternative to `Exit()` in Cobra commands. See #5.
		}
		sh.AddDispatcher(func(args []string) (bool, error) {
			rootCmd.SetArgs(args)
			rootCmd.Use = "."
			err := rootCmd.Execute()
			if err != nil {
				if strings.HasPrefix(err.Error(), "unknown command") {
					return false, nil
				}
				return false, err
			}
			return true, nil
		})

		sh.Run()
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
