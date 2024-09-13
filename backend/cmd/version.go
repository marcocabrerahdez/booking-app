package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Backend",
	Long:  `Print the version number of Backend`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Backend CLI v0.0.1")
		fmt.Println("Config file:", viper.ConfigFileUsed())
	},
}
