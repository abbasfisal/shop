package customer

import (
	"shop/domain/entities"
	"shop/domain/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) repositories.CustomerRepositoryInterface {
	return &CustomerRepository{db: db}
}

func (cr *CustomerRepository) GetAll(c *gin.Context) ([]*entities.Customer, error) {

	var customers []*entities.Customer
	if err := cr.db.WithContext(c).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}
