package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"shop/internal/database/mysql"
	"shop/internal/entities"
	adminAuthRepo "shop/internal/modules/admin/repositories/auth"
	"shop/internal/pkg/sessions"
	"strconv"
)

func Auth(c *gin.Context) entities.User {
	authID := sessions.GET(c, "auth_id")
	fmt.Println("authid:", authID)
	if authID == "" {
		return entities.User{}
	}

	repo := adminAuthRepo.NewAuthenticateRepository(mysql.Get())
	userID, _ := strconv.Atoi(authID)
	user, _ := repo.FindByUserID(c, uint(userID))

	return user
}
