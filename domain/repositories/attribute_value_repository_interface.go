package repositories

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"
)

type AttributeValueRepositoryInterface interface {
	// Store builds the entity from the admin form (color hex lands in Meta
	// only when the parent attribute is a color attribute).
	Store(ctx context.Context, req *requests.CreateAttributeValueRequest) (*entities.AttributeValue, error)
	GetAllAttribute(c *gin.Context) ([]*entities.Attribute, error)
	Find(c *gin.Context, attributeValueID int) (*entities.AttributeValue, error)
	Update(c *gin.Context, attributeValueID int, req *requests.UpdateAttributeValueRequest) (*entities.AttributeValue, error)
}
