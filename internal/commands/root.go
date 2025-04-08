package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"shop/internal/pkg/bootstrap"
)

var rootCmd = &cobra.Command{
	Use:   "help",
	Short: "display help ",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {

	// load configs to access .env keys in commands
	err := bootstrap.LoadConfig()
	if err != nil {
		log.Fatalln("failed load config : ", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
