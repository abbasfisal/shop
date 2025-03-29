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
	"strings"
)

func (a *AdminHandler) CreateAttribute(c *gin.Context) {
	html.Render(c, http.StatusFound, "admin_create_attribute",
		gin.H{
			"TITLE": "ایجاد اتریبیوت",
		})
	return
}

func (a *AdminHandler) StoreAttribute(c *gin.Context) {
	//todo: check uniqueness of title in given category

	var req requests.CreateAttributeRequest

	_ = c.Request.ParseForm()
	err := c.ShouldBind(&req)
	if err != nil {
		fmt.Println("err:", err)
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/attributes/create")
		return
	}

	newAttr, Aerr := a.attributeSrv.Create(c, &req)

	if Aerr != nil || newAttr.ID <= 0 {
		fmt.Println("Error in creating new attribute : ", err)
		sessions.Set(c, "message", custom_messages.AttributeCreateFailed)
		c.Redirect(http.StatusFound, "/admins/attributes/create")
		return
	}

	sessions.Set(c, "message", custom_messages.AttributeCreateSuccessful)
	c.Redirect(http.StatusFound, "/admins/attributes/create")
	return
}

func (a *AdminHandler) IndexAttribute(c *gin.Context) {

	attributes, err := a.attributeSrv.Index(c)

	if err.Code == 400 {
		c.JSON(200, gin.H{
			"data": "empty",
		})
	} else if err.Code == 500 {
		html.Error500(c)
		return
	}
	html.Render(c, 200, "admin_index_attribute",
		gin.H{
			"TITLE":      "لیست اتریبیوت ها",
			"ATTRIBUTES": attributes,
		})
	return

}

func (a *AdminHandler) ShowAttribute(c *gin.Context) {
	attributeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}

	attributeShow, attErr := a.attributeSrv.Show(c, attributeID)
	if attErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}
	if attErr.Code == 500 {
		html.Error500(c)
		return
	}

	html.Render(c, http.StatusFound, "admin_show_attribute",
		gin.H{
			"TITLE":     "نمایش اتریبیوت",
			"ATTRIBUTE": attributeShow,
		},
	)
}

func (a *AdminHandler) UpdateAttribute(c *gin.Context) {

	attributeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}

	url := fmt.Sprintf("/admins/attributes/%d/edit", attributeID)

	oldAttribute, oldErr := a.attributeSrv.Show(c, attributeID)

	if oldErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}
	if oldErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}

	var req requests.CreateAttributeRequest

	_ = c.Request.ParseForm()
	bErr := c.ShouldBind(&req)
	if bErr != nil {
		errors.Init()
		errors.SetFromErrors(bErr)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, url)
		return
	}

	//don't need to update
	if oldAttribute.Title == strings.TrimSpace(req.Title) {
		sessions.Set(c, "message", custom_messages.AttributeUpdatedSuccessfully)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}

	updateErr := a.attributeSrv.Update(c, attributeID, &req)
	if updateErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}
	if updateErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}

	sessions.Set(c, "message", custom_messages.AttributeUpdatedSuccessfully)
	c.Redirect(http.StatusFound, "/admins/attributes")
}

func (a *AdminHandler) AppendAttribute(c *gin.Context) {

	//convert string id to int
	inventoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	attributes := c.PostFormArray("attributes")
	if len(attributes) <= 0 {
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	//store attributes
	pAttrErr := a.productSrv.AppendAttributesToInventory(c, inventoryID, attributes)
	if pAttrErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}
	if pAttrErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	sessions.Set(c, "message", custom_error.SuccessfullyCreated)

	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}

func (a *AdminHandler) GetAttributesByCategoryID(c *gin.Context) {
	//todo: error for converting string to integer
	cat, err := strconv.Atoi(c.Param("catID"))
	fmt.Println(" category id : ", cat)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	attributes, err := a.attributeSrv.FetchByCategoryID(c, cat)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attributes"})
		return
	}

	fmt.Println("response attributes : ", attributes)
	c.JSON(http.StatusOK, attributes)
}
