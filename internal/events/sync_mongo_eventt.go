package events

import (
	"context"
	"fmt"
	"shop/internal/modules/admin/repositories/product"
)

const SyncMongoEvent = "sync.mongo"

type SyncMongoEventPayload struct {
	ProductID uint
}

func SyncMongoListener(ctx context.Context, data any) {
	payload := data.(SyncMongoEventPayload)
	dep := GetDep()

	select {
	case <-ctx.Done():
		fmt.Println("SyncMongo : Execution canceled or timed out")
		return
	default:
		product.SyncMongo(ctx, dep.DB, payload.ProductID) //todo: call Asynq to update
	}
}
