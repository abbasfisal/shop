package attributeValue

import (
	"context"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type AttributeValueServiceInterface interface {
	Create(ctx context.Context, req requests.CreateAttributeValueRequest) (responses.AttributeValue, custom_error.CustomError)
}
