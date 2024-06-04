package commands

import (
	"github.com/spf13/cobra"
	"shop/internal/database/mysql/seeder"
)

func init() {
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed  Tables",
	Run: func(cmd *cobra.Command, args []string) {
		seeder.Seed()
	},
}
