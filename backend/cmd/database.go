package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(databaseCmd)
}

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database commands",
	Long:  `Database commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
