package repositories

import (
	"context"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"

	"github.com/gin-gonic/gin"
)

type AttributeRepositoryInterface interface {
	Store(ctx context.Context, attr *entities.Attribute) (*entities.Attribute, error)
	GetByCategory(ctx context.Context, catID int) ([]*entities.Attribute, error)
	GetAll(c *gin.Context) ([]*entities.Attribute, error)
	GetByID(c context.Context, attributeID int) (*entities.Attribute, error)
	Update(c *gin.Context, attributeID int, req *requests.CreateAttributeRequest) error
}
