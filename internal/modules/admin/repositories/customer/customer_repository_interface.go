package customer

import (
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
)

type CustomerRepositoryInterface interface {
	GetAll(c *gin.Context) ([]entities.Customer, error)
}
