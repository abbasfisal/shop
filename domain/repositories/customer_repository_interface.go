package repositories

import (
	"github.com/gin-gonic/gin"
	"shop/domain/entities"
)

type CustomerRepositoryInterface interface {
	GetAll(c *gin.Context) ([]*entities.Customer, error)
}
