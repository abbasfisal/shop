package sessions

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessions", store))
}

func Set(c *gin.Context, key, value string) {
	session := sessions.Default(c)

	session.Set(key, value)
	session.Save()
}

func Flash(c *gin.Context, key string) string {
	session := sessions.Default(c)

	response := session.Get(key)
	session.Save()

	session.Delete(key)
	session.Save()

	if response != nil {
		return response.(string)
	}
	return ""
}

func GET(c *gin.Context, key string) string {
	session := sessions.Default(c)

	response := session.Get(key)
	session.Save()

	if response != nil {
		return response.(string)
	}
	return ""
}

func Remove(c *gin.Context, key string) {
	session := sessions.Default(c)

	session.Delete(key)
	session.Save()
}

func ClearAll(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}
