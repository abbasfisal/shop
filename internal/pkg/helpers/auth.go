package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"shop/internal/database/mysql"
	"shop/internal/entities"
	adminAuthRepo "shop/internal/modules/admin/repositories/auth"
	customerAuthRepo "shop/internal/modules/public/repositories/customer_auth"
	customerResponse "shop/internal/modules/public/responses"
	"shop/internal/pkg/sessions"
	"strconv"
)

func Auth(c *gin.Context) entities.User {
	authID := sessions.GET(c, "auth_id")
	fmt.Println("helpers - Auth() - authid:", authID)
	if authID == "" {
		return entities.User{}
	}

	repo := adminAuthRepo.NewAuthenticateRepository(mysql.Get())
	userID, _ := strconv.Atoi(authID)
	user, _ := repo.FindByUserID(c, uint(userID))

	return user
}

func CustomerAuth(c *gin.Context) customerResponse.Customer {
	sessionID := sessions.GET(c, "session_id") //session uuid

	if sessionID == "" {
		return customerResponse.Customer{}
	}

	repo := customerAuthRepo.NewAuthenticateRepository(mysql.Get())
	customer, err := repo.FindCustomerBySessionID(c, sessionID)
	if err != nil {
		return customerResponse.Customer{}
	}

	return customerResponse.ToCustomer(customer)
}
