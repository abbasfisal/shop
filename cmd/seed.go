package cmd

import (
	"github.com/spf13/cobra"
	"shop/infrastructure/database/postgres"
	"shop/infrastructure/seeders"
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
