package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/internal/modules/admin/requests"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/custom_messages"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/old"
	"shop/internal/pkg/sessions"
	"strconv"
)

func (a *AdminHandler) CreateProductFeature(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	productShow, _, pErr := a.productSrv.Show(c, "id", pID)

	if pErr.Code > 0 {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	html.Render(c, http.StatusFound, "admin_add_product_feature",
		gin.H{
			"TITLE":   "افزودن feature به محصول",
			"PRODUCT": productShow,
		},
	)
	return
}

func (a *AdminHandler) StoreProductFeature(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	url := fmt.Sprintf("/admins/products/%d/add-feature", pID)
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	_, _, pErr := a.productSrv.Show(c, "id", pID)
	if pErr.Code > 0 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	var req requests.CreateProductFeatureRequest
	_ = c.Request.ParseForm()
	if bErr := c.ShouldBind(&req); bErr != nil {
		fmt.Println("---  product feature bind Err : ", bErr)
		errors.Init()
		errors.SetFromErrors(bErr)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, url)
		return
	}

	addErr := a.productSrv.AddFeature(c, pID, &req)
	if addErr.Code > 0 {
		fmt.Println("-- add feature error :", addErr)
		sessions.Set(c, "errors", custom_error.SomethingWrongHappened)
		c.Redirect(http.StatusFound, url)
		return
	}

	sessions.Set(c, "message", custom_error.SuccessfullyCreated)
	c.Redirect(http.StatusFound, url)
	return
}

func (a *AdminHandler) ShowProductFeature(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	//url := fmt.Sprintf("/admins/products/%d/show-feature", pID)
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	productData, _, pErr := a.productSrv.Show(c, "id", pID)
	if pErr.Code > 0 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	html.Render(c, http.StatusOK, "admin_show_product_feature", gin.H{
		"TITLE":   "",
		"PRODUCT": productData,
	})
	return
}

func (a *AdminHandler) DeleteProductFeature(c *gin.Context) {
	pID, pErr := strconv.Atoi(c.Param("id"))
	fID, fErr := strconv.Atoi(c.Param("featureID"))
	url := fmt.Sprintf("/admins/products/%d/show-feature", pID)
	if pErr != nil || fErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, url)
		return
	}

	deleteErr := a.productSrv.RemoveFeature(c, pID, fID)
	if deleteErr.Code > 0 {
		fmt.Println("--- delete feature err  : ", deleteErr)
		sessions.Set(c, "message", custom_error.SomethingWrongHappened)
		c.Redirect(http.StatusFound, url)
	}

	sessions.Set(c, "message", custom_messages.DeleteSuccessfully)
	c.Redirect(http.StatusFound, url)

}

func (a *AdminHandler) EditProductFeature(c *gin.Context) {
	pID, pErr := strconv.Atoi(c.Param("id"))
	fID, fErr := strconv.Atoi(c.Param("featureID"))
	url := fmt.Sprintf("/admins/products/%d/show-feature", pID)
	if pErr != nil || fErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, url)
		return
	}

	feat, err := a.productSrv.FetchFeature(c, pID, fID)
	if err.Code > 0 {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, url)
		return
	}

	html.Render(c, http.StatusFound, "admin_edit_product_feature", gin.H{
		"TITLE":   "ویرایش صفت",
		"FEATURE": feat,
	})

}

func (a *AdminHandler) UpdateProductFeature(c *gin.Context) {
	pID, pErr := strconv.Atoi(c.Param("id"))
	fID, fErr := strconv.Atoi(c.Param("featureID"))
	url := fmt.Sprintf("/admins/products/%d/show-feature", pID)
	if pErr != nil || fErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, url)
		return
	}

	var req requests.UpdateProductFeatureRequest
	_ = c.Request.ParseForm()
	if bErr := c.ShouldBind(&req); bErr != nil {
		fmt.Println("---  product feature bind Err : ", bErr)
		errors.Init()
		errors.SetFromErrors(bErr)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, url)
		return
	}

	updateErr := a.productSrv.UpdateFeature(c, pID, fID, &req)
	if updateErr.Code > 0 {
		fmt.Println("-- update feature error :", updateErr)
		sessions.Set(c, "errors", custom_error.SomethingWrongHappened)
		c.Redirect(http.StatusFound, url)
		return
	}

	sessions.Set(c, "message", custom_error.SuccessfullyCreated)
	c.Redirect(http.StatusFound, url)
	return

}
