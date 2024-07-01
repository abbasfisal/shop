package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/internal/database/mysql"
	adminAuthRepo "shop/internal/modules/admin/repositories/auth"
	"shop/internal/pkg/sessions"
	"strconv"
)

func IsAdmin(c *gin.Context) {
	fmt.Println("admin middleware")

	authID := sessions.GET(c, "auth_id")

	if authID == "" {
		fmt.Println("auth id not found in admin middleware")
		c.Redirect(http.StatusFound, "/admins/login")
		return
	}

	repo := adminAuthRepo.NewAuthenticateRepository(mysql.Get())

	userID, _ := strconv.Atoi(authID)
	user, _ := repo.FindByUserID(c, uint(userID))

	if user.ID <= 0 || user.Type != "admin" {
		c.Redirect(http.StatusFound, "/admins/login")
	}

	c.Next()
}
