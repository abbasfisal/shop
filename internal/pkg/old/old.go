package old

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

var oldList = make(map[string][]string)

func Init() {
	oldList = map[string][]string{}
}

func Set(c *gin.Context) {
	c.Request.ParseForm()
	oldList = c.Request.PostForm
}

func Get() map[string][]string {
	return oldList
}

func ToString() string {
	out, _ := json.Marshal(Get())
	if out != nil {
		return string(out)
	}
	return ""
}
