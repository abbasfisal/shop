package attributeValue

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type AttributeValueServiceInterface interface {
	Create(ctx context.Context, req requests.CreateAttributeValueRequest) (responses.AttributeValue, custom_error.CustomError)
	IndexAttribute(c *gin.Context) (responses.Attributes, custom_error.CustomError)
	Show(c *gin.Context, attributeValueID int) (responses.AttributeValue, custom_error.CustomError)
}
