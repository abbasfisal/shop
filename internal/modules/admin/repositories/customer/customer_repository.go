package customer

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return CustomerRepository{db: db}
}

func (cr CustomerRepository) GetAll(c *gin.Context) ([]entities.Customer, error) {

	var customers []entities.Customer
	if err := cr.db.WithContext(c).Find(&customers).Error; err != nil {
		return customers, err
	}
	return customers, nil
}
