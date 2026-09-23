package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/infrastructure/database/postgres"
	adminAuthRepo "shop/infrastructure/repositories/auth"
	"shop/pkg/sessions"
	"strconv"
)

func IsGuest(c *gin.Context) {

	fmt.Println("Guest middleware ")

	authID := sessions.GET(c, "auth_id")
	fmt.Println("middleware - authid:", authID)
	if authID == "" {
		fmt.Println("guest auth not found")
		sessions.Remove(c, "auth_id")
		c.Next()
		return
	}

	repo := adminAuthRepo.NewAuthenticateRepository(postgres.Get())
	userID, _ := strconv.Atoi(authID)
	user, _ := repo.FindByUserID(c, uint(userID))

	//user was find
	if user.ID > 0 && user.Type == "admin" {
		fmt.Println("user is logged in and is admin")
		c.Redirect(http.StatusFound, "/admins/home")
		return
	}
	c.Next()
}
