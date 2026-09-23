package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/typesense/typesense-go/v3/typesense/api"
	"github.com/typesense/typesense-go/v3/typesense/api/pointer"
	"shop/interfaces/http/response"
)

// SearchProductByTypesence realtime product search.
// Response shape is fixed by tsearch.html: {results:[{title, link}]}.
func (p PublicHandler) SearchProductByTypesence(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	if query == "" {
		c.JSON(http.StatusOK, gin.H{"results": []string{}})
		return
	}

	if p.dep.TypeSenceClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "جستجو در حال حاضر در دسترس نیست"})
		return
	}

	searchParams := &api.SearchCollectionParams{
		Q:       pointer.String(query),
		QueryBy: pointer.String("title,sku,description,category,brand"),
	}
	searchResults, err := p.dep.TypeSenceClient.
		Collection("products").Documents().
		Search(c.Request.Context(), searchParams)

	if err != nil && strings.Contains(err.Error(), "Could not infer field") {
		// older collection schema without the rich fields — retry on the
		// legacy field set so search keeps working until search:reindex --recreate
		searchParams.QueryBy = pointer.String("title,sku")
		searchResults, err = p.dep.TypeSenceClient.
			Collection("products").Documents().
			Search(c.Request.Context(), searchParams)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در جستجو"})
		return
	}

	results := []map[string]interface{}{}
	if searchResults != nil && searchResults.Hits != nil {
		for _, hit := range *searchResults.Hits {
			if hit.Document == nil {
				continue
			}
			document := *hit.Document

			title, _ := document["title"].(string)
			sku, _ := document["sku"].(string)
			slug, _ := document["slug"].(string)
			link := fmt.Sprintf("http://%s:%s/product/%s/%s",
				os.Getenv("APP_HOST"),
				os.Getenv("APP_PORT"),
				sku,
				slug,
			)

			results = append(results, map[string]interface{}{
				"title": title,
				"link":  link,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (p PublicHandler) ShowTypeSenceForm(c *gin.Context) {
	response.CustomerRender(c, http.StatusFound, "tsearch",
		gin.H{
			"title": "search",
		})
	return
}
