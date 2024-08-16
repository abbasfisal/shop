package customer_auth

import (
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
)

type AuthenticateRepositoryInterface interface {
	// FindCustomerBySessionID : sessionID is uuid
	FindCustomerBySessionID(c *gin.Context, sessionID string) (entities.Customer, error)
}
