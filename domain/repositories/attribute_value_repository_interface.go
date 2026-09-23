package repositories

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"
)

type AttributeValueRepositoryInterface interface {
	Store(ctx context.Context, attr *entities.AttributeValue) (*entities.AttributeValue, error)
	GetAllAttribute(c *gin.Context) ([]*entities.Attribute, error)
	Find(c *gin.Context, attributeValueID int) (*entities.AttributeValue, error)
	Update(c *gin.Context, attributeValueID int, req *requests.UpdateAttributeValueRequest) (*entities.AttributeValue, error)
}
