package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/typesense/typesense-go/v3/typesense/api"
	"github.com/typesense/typesense-go/v3/typesense/api/pointer"
	"log"
	"net/http"
	"os"
	"shop/internal/pkg/html"
)

func (p PublicHandler) SearchProductByTypesence(c *gin.Context) {

	query := c.DefaultQuery("q", "")
	if query == "" {
		log.Println("0")
		c.JSON(http.StatusOK, gin.H{"results": []string{}})
		return
	}

	log.Println("1")

	// search in Typesense
	searchParams := &api.SearchCollectionParams{
		Q:       pointer.String(query),
		QueryBy: pointer.String("title,sku"),
	}
	searchResults, err := p.dep.TypeSenceClient.
		Collection("products").Documents().
		Search(c.Request.Context(), searchParams)

	if err != nil {
		log.Println(2)
		log.Println("search error :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در جستجو"})
		return
	}
	log.Println(3)

	//var results []string
	var results []map[string]interface{}
	log.Println(4)

	for _, hit := range *searchResults.Hits {
		log.Println(5)

		document := *hit.Document

		title := document["title"].(string)
		link := fmt.Sprintf("http://%s:%s/product/%s/%s",
			os.Getenv("APP_HOST"),
			os.Getenv("APP_PORT"),
			document["sku"].(string),
			document["slug"].(string),
		)

		results = append(results, map[string]interface{}{
			"title": title,
			"link":  link,
		})

	}
	log.Println(6)

	c.JSON(http.StatusOK, gin.H{"results": results})

}

func (p PublicHandler) ShowTypeSenceForm(c *gin.Context) {

	search, err := p.dep.TypeSenceClient.Collection("products").Documents().Search(c.Request.Context(), &api.SearchCollectionParams{
		Q: pointer.String("*"),
	})
	c.JSON(200, gin.H{
		"err":  err,
		"data": search,
	})

	return

	html.CustomerRender(c, http.StatusFound, "tsearch",
		gin.H{
			"title": "search",
		},
	)
	return
}
