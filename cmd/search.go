package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"shop/bootstrap"
	"shop/infrastructure/database/postgres"
	"shop/infrastructure/database/typesenceclient"
	"shop/infrastructure/repositories/product"
)

var reindexRecreate bool

func init() {
	rootCmd.AddCommand(reindexCmd)
	reindexCmd.Flags().BoolVar(&reindexRecreate, "recreate", false,
		"drop and recreate the typesense collection before reindexing (required after schema changes)")
}

var reindexCmd = &cobra.Command{
	Use:   "search:reindex",
	Short: "Rebuild the typesense product index from PostgreSQL read models",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := bootstrap.Initialize(); err != nil {
			log.Fatal("[x] error initializing project: ", err)
		}
		postgres.Connect()

		client := typesenceclient.GetTClient()
		if client == nil {
			log.Fatal("[x] typesense client not available (check TYPESENCE_HOST/PORT/API_KEY)")
		}

		if reindexRecreate {
			if err := typesenceclient.RecreateSchema(client); err != nil {
				log.Fatal("[x] recreate schema failed: ", err)
			}
			fmt.Println("typesense schema recreated")
		}

		db := postgres.Get()
		var ids []uint
		if err := db.Table("products").Order("id").Pluck("id", &ids).Error; err != nil {
			log.Fatal("[x] load products failed: ", err)
		}

		ctx := context.Background()
		ok, failed := 0, 0
		for _, id := range ids {
			if err := product.SyncReadModel(ctx, db, id); err != nil {
				failed++
				log.Printf("[reindex] product %d failed: %v", id, err)
				continue
			}
			ok++
		}

		fmt.Printf("reindex done: %d indexed, %d failed (of %d)\n", ok, failed, len(ids))
	},
}
