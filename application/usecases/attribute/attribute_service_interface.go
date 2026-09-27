package attribute

import (
	"context"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/interfaces/http/requests/admin"

	"github.com/gin-gonic/gin"
)

type AttributeServiceInterface interface {
	Create(ctx context.Context, req *requests.CreateAttributeRequest) (*responses.Attribute, error)
	FetchByCategoryID(ctx context.Context, categoryID int) (*responses.Attributes, error)
	Index(c *gin.Context) (*responses.Attributes, domain_err.CustomError)
	Show(c context.Context, attributeID int) (*responses.Attribute, domain_err.CustomError)
	Update(c *gin.Context, attributeID int, req *requests.CreateAttributeRequest) domain_err.CustomError
}
