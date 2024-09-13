package cmd

import (
	"backend/api"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the Backend server",
	Long:  `Starts the Backend server, usefull for development.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting %s server...", viper.GetString("general.app.name"))
		api.Serve()
	},
}
