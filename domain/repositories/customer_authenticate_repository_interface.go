package repositories

import (
	"shop/domain/entities"

	"github.com/gin-gonic/gin"
)

type CustomerAuthenticateRepositoryInterface interface {
	// FindCustomerBySessionID : sessionID is uuid
	FindCustomerBySessionID(c *gin.Context, sessionID string) (entities.Customer, error)
}
