package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"zdora/config"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show zdora config",
	Long: `Show zdora config
Should be config first!
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			params := strings.Join(args, " ")
			config.SetConfig(params)
		} else {
			fmt.Println(config.GetConfigSchema())
		}
	},
}
