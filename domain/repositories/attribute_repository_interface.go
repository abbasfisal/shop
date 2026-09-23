package repositories

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
	"shop/interfaces/http/requests/admin"
)

type AttributeRepositoryInterface interface {
	Store(ctx context.Context, attr *entities.Attribute) (*entities.Attribute, error)
	GetByCategory(ctx context.Context, catID int) ([]*entities.Attribute, error)
	GetAll(c *gin.Context) ([]*entities.Attribute, error)
	GetByID(c context.Context, attributeID int) (*entities.Attribute, error)
	Update(c *gin.Context, attributeID int, req *requests.CreateAttributeRequest) error
}
