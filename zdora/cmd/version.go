package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of zdora",
	Long:  `All software has versions. This is zdora's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("zdora version is v1.0")
	},
}
