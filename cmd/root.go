package cmd

import (
	"os"

	core "github.com/seabreeze-project/seabreeze/core"
	"github.com/seabreeze-project/seabreeze/util"
	"github.com/spf13/cobra"
)

var (
	Core = core.New()

	// flags
	configFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "seabreeze",
	Short: "Seabreeze",
	Long:  `Seabreeze: A really simple container orchestration tool with superpowers`,
}

var h = util.NewHelper(rootCmd)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "use given config file")
}

func initConfig() {
	if err := Core.LoadConfig(configFile); err != nil {
		h.Fatal("cannot load config:", err)
	}

	if err := Core.Config().Validate(); err != nil {
		h.Fatal("error validating config:", err)
	}
}
