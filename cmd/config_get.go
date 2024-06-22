package cmd

import (
	"github.com/seabreeze-project/seabreeze/console"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// configGetCmd represents the 'config get' command
var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a configuration value",
	Long:  `This command returns the value of the given configuration key.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		h := console.NewHelper(cmd)

		val := Core.Config().Get(args[0])
		if val == nil {
			h.Fatalf("Key %q undefined\n", args[0])
		}
		out, err := yaml.Marshal(val)
		if err != nil {
			h.Fatal(err)
		}
		h.Println(string(out))
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
}
