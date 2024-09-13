package cmd

import (
	"backend/database"

	"github.com/spf13/cobra"
)

func init() {
	databaseCmd.AddCommand(dropCmd)
}

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drops the database",
	Long:  `Drops the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := database.Connect(); err != nil {
			panic(err)
		}
		if err := database.Drop(); err != nil {
			panic(err)
		}
	},
}
