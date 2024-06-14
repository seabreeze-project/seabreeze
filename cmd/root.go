package cmd

import (
	"fmt"
	"os"

	core "github.com/seabreeze-project/seabreeze/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	if configFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(Core.ConfigBasePath)
		viper.SetConfigType("yml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("SEABREEZE")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
