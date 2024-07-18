package attribute

import (
	"context"
	"shop/internal/entities"
)

type AttributeRepositoryInterface interface {
	Store(ctx context.Context, attr entities.Attribute) (entities.Attribute, error)
}
