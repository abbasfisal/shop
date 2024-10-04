package customer_auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
)

type AuthenticateRepository struct {
	db *gorm.DB
}

func NewAuthenticateRepository(db *gorm.DB) AuthenticateRepository {
	return AuthenticateRepository{
		db: db,
	}
}

func (ar AuthenticateRepository) FindCustomerBySessionID(c *gin.Context, sessionID string) (entities.Customer, error) {
	var sess entities.Session

	err := ar.db.
		WithContext(c).
		Preload("Customer.Address").
		Preload("Customer.Carts.CartItems").
		Where("session_id = ?", sessionID).
		First(&sess).
		Error

	if err != nil {
		return entities.Customer{}, err
	}

	return sess.Customer, nil
}
