package repositories

import (
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
)

type CustomerAuthenticateRepositoryInterface interface {
	// FindCustomerBySessionID : sessionID is uuid
	FindCustomerBySessionID(c *gin.Context, sessionID string) (entities.Customer, error)
}
