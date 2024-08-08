package attributeValue

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
)

type AttributeValueRepositoryInterface interface {
	Store(ctx context.Context, attr entities.AttributeValue) (entities.AttributeValue, error)
	GetAllAttribute(c *gin.Context) ([]entities.Attribute, error)
}
