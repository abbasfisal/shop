package cmd

import (
	"shop/infrastructure/database/postgres"
	"shop/infrastructure/seeders"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database tables",
	Run: func(cmd *cobra.Command, args []string) {
		postgres.Connect()
		seeders.Seed()
	},
}
