package events

import (
	"context"
	"fmt"
	"shop/infrastructure/repositories/product"
)

const SyncReadModelEvent = "sync.readmodel"

type SyncReadModelEventPayload struct {
	ProductID uint
}

func SyncReadModelListener(ctx context.Context, data any) {
	payload := data.(SyncReadModelEventPayload)
	dep := GetDep()

	select {
	case <-ctx.Done():
		fmt.Println("SyncReadModel : Execution canceled or timed out")
		return
	default:
		product.SyncReadModel(ctx, dep.DB, payload.ProductID)
	}
}
