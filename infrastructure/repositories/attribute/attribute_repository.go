package attribute

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
	"strings"
)

type AttributeRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) repositories.AttributeRepositoryInterface {
	return &AttributeRepository{db: db}
}

//------------------
//>>>>> Methods <<<<
//------------------

func (ar *AttributeRepository) Store(ctx context.Context, attr *entities.Attribute) (*entities.Attribute, error) {
	err := ar.db.WithContext(ctx).Create(&attr).Error
	return attr, err
}

// GetByCategory is kept for the legacy AJAX route
// (/admins/get-attributes/:catID): attributes are no longer category scoped,
// so it returns every attribute with its values (sorted like the Laravel model).
func (ar *AttributeRepository) GetByCategory(ctx context.Context, catID int) ([]*entities.Attribute, error) {
	return ar.getAll(ctx)
}
func (ar *AttributeRepository) GetAll(c *gin.Context) ([]*entities.Attribute, error) {
	return ar.getAll(c)
}

// getAll loads attributes with values (ordered by sort_order) — required by
// the product create/edit combination builder.
func (ar *AttributeRepository) getAll(ctx context.Context) ([]*entities.Attribute, error) {
	var attributes []*entities.Attribute
	err := ar.db.WithContext(ctx).
		Preload("AttributeValues", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		}).
		Order("sort_order ASC, id ASC").
		Find(&attributes).Error
	return attributes, err
}
func (ar *AttributeRepository) GetByID(c context.Context, attributeID int) (*entities.Attribute, error) {
	var att entities.Attribute
	err := ar.db.WithContext(c).Preload("AttributeValues").First(&att, attributeID).Error
	return &att, err
}
func (ar *AttributeRepository) Update(c *gin.Context, attributeID int, req *requests.CreateAttributeRequest) error {
	var att entities.Attribute
	err := ar.db.First(&att, attributeID).Error
	if err != nil {
		return err
	}

	return ar.db.Model(&att).Updates(map[string]interface{}{
		"title":      strings.TrimSpace(req.Title),
		"input_type": req.NormalizedInputType(),
	}).Error
}
