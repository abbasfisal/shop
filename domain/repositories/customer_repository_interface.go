package repositories

import (
	"shop/domain/entities"

	"github.com/gin-gonic/gin"
)

type CustomerRepositoryInterface interface {
	GetAll(c *gin.Context) ([]*entities.Customer, error)
}
