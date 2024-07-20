package attribute

import (
	"context"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
)

type AttributeServiceInterface interface {
	Create(ctx context.Context, req requests.CreateAttributeRequest) (responses.Attribute, error)
	FetchByCategoryID(ctx context.Context, categoryID int) (responses.Attributes, error)
}
