package attribute

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type AttributeServiceInterface interface {
	Create(ctx context.Context, req requests.CreateAttributeRequest) (responses.Attribute, error)
	FetchByCategoryID(ctx context.Context, categoryID int) (responses.Attributes, error)
	Index(c *gin.Context) (responses.Attributes, custom_error.CustomError)
	Show(c context.Context, attributeID int) (responses.Attribute, custom_error.CustomError)
	Update(c *gin.Context, attributeID int, req requests.CreateAttributeRequest) custom_error.CustomError
}
