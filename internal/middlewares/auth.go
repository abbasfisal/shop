package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/internal/database/mysql"
	adminAuthRepo "shop/internal/modules/admin/repositories/auth"
	"shop/internal/pkg/sessions"
	"strconv"
)

func IsAdmin(c *gin.Context) {

	authID := sessions.GET(c, "auth_id")
	if authID == "" {
		c.Redirect(http.StatusFound, "/admins/login")
		return
	}

	userID, err := strconv.Atoi(authID)
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/login")
		return
	}

	repo := adminAuthRepo.NewAuthenticateRepository(mysql.Get())

	user, err := repo.FindByUserID(c, uint(userID))

	if err != nil || user == nil || user.ID <= 0 || user.Type != "admin" {
		c.Redirect(http.StatusFound, "/admins/login")
		return
	}

	c.Next()
}
