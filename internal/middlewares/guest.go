package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	adminAuthRepo "shop/internal/modules/admin/repositories/auth"
	"shop/internal/pkg/sessions"
	"strconv"
)

func IsGuest(c *gin.Context) {

	fmt.Println("Guest middleware ")

	authID := sessions.GET(c, "auth_id")
	fmt.Println("authid:", authID)
	if authID == "" {
		fmt.Println("guest auth not found")
		sessions.Remove(c, "auth_id")
		c.Next()
		return
	}

	repo := adminAuthRepo.NewAuthenticateRepository()
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
