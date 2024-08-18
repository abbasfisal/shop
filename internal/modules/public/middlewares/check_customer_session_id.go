package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/internal/pkg/helpers"
)

func CheckCustomerSessionID() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 1. دسترسی به پروتکل (http یا https)
		scheme := c.Request.URL.Scheme
		if scheme == "" {
			scheme = "http"
			if c.Request.TLS != nil {
				scheme = "https"
			}
		}

		// 2. دسترسی به هاست (مثلاً example.com)
		host := c.Request.Host

		// 3. دسترسی به مسیر (path)
		path := c.Request.URL.Path

		// 4. دسترسی به کوئری استرینگ کامل
		rawQuery := c.Request.URL.RawQuery

		// 5. دسترسی به URL کامل (پروتکل + هاست + مسیر + کوئری)
		fullURL := scheme + "://" + host + path
		if rawQuery != "" {
			fullURL += "?" + rawQuery
		}

		// 6. دسترسی به کوئری پارامترهای جداگانه
		id := c.Query("id")
		name := c.Query("name")

		// 7. دسترسی به متد HTTP (GET, POST و ...)
		method := c.Request.Method
		//
		//// چاپ اطلاعات
		//fmt.Println("Scheme:", scheme)
		//fmt.Println("Host:", host)
		//fmt.Println("Path:", path)
		//fmt.Println("Raw Query:", rawQuery)
		//fmt.Println("Full URL:", fullURL)
		//fmt.Println("Query Parameter ID:", id)
		//fmt.Println("Query Parameter Name:", name)
		//fmt.Println("HTTP Method:", method)

		// ارسال پاسخ به کلاینت
		fmt.Println(" ~~~ URL DATA : ", gin.H{
			"scheme":      scheme,
			"host":        host,
			"path":        path,
			"raw_query":   rawQuery,
			"full_url":    fullURL,
			"query_id":    id,
			"query_name":  name,
			"http_method": method,
		})

		fmt.Println("Middleware : check_customer_session")

		customer := helpers.CustomerAuth(c)

		if customer.ID > 0 {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}
