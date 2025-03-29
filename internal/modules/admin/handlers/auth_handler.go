package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/internal/modules/admin/requests"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/old"
	"shop/internal/pkg/sessions"
	"strconv"
)

func (a *AdminHandler) ShowLogin(c *gin.Context) {
	html.Render(c, http.StatusOK, "modules/admin/html/admin_login", gin.H{"title": "login"})
	return
}

func (a *AdminHandler) PostLogin(c *gin.Context) {
	var req requests.LoginRequest
	_ = c.Request.ParseForm()

	if err := c.ShouldBind(&req); err != nil {
		errors.SetErrors(c, a.dep.I18nBundle, err)

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/login")
		return
	}

	user, loginErr := a.authSrv.Login(c.Request.Context(), &req)
	if loginErr.Error() != "" {
		if loginErr.Code == 404 {
			html.Render(c, http.StatusFound, "modules/admin/html/admin_login", gin.H{
				"MESSAGE": loginErr.Error(),
			})
			return
		}
		if loginErr.Code == 500 {
			html.Error500(c)
			return
		}
	}

	sessions.Set(c, "auth_id", strconv.Itoa(int(user.ID)))
	c.Redirect(http.StatusFound, "/admins/home")
}
