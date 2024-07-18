package attributeValue

import (
	"context"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
)

type AttributeValueServiceInterface interface {
	Create(ctx context.Context, req requests.CreateAttributeValueRequest) (responses.AttributeValue, error)
}
