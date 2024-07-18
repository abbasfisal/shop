package attributeValue

import (
	"context"
	"shop/internal/entities"
)

type AttributeValueRepositoryInterface interface {
	Store(ctx context.Context, attr entities.AttributeValue) (entities.AttributeValue, error)
}
