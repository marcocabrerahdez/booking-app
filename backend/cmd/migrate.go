package cmd

import (
	"backend/database"

	"github.com/spf13/cobra"
)

func init() {
	databaseCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates the database",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := database.Connect(); err != nil {
			panic(err)
		}
		if err := database.Migrate(); err != nil {
			panic(err)
		}
	},
}
