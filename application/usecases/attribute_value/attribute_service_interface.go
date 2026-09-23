package attributeValue

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/interfaces/http/requests/admin"
)

type AttributeValueServiceInterface interface {
	Create(ctx context.Context, req *requests.CreateAttributeValueRequest) (*responses.AttributeValue, domain_err.CustomError)
	IndexAttribute(c *gin.Context) (*responses.Attributes, domain_err.CustomError)
	Show(c *gin.Context, attributeValueID int) (*responses.AttributeValue, domain_err.CustomError)
	Update(c *gin.Context, attributeValueID int, req *requests.UpdateAttributeValueRequest) domain_err.CustomError
}
