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

func (a *AdminHandler) CreateAttributeValues(c *gin.Context) {

	attributes, _ := a.attributeSrv.Index(c)
	html.Render(c, http.StatusFound, "admin_create_attribute_values", gin.H{
		"TITLE":      "create new attribute-values",
		"ATTRIBUTES": attributes,
	})
	return
}

func (a *AdminHandler) StoreAttributeValues(c *gin.Context) {

	var req requests.CreateAttributeValueRequest
	_ = c.Request.ParseForm()

	err := c.ShouldBind(&req)
	if err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		if req.AttributeID <= 0 {
			errors.Add("attribute_id", custom_messages.SelectOne)
		}
		fmt.Println("validation failed : ", err.Error())
		fmt.Println("attribute id : ", req.AttributeID)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/attribute-values/create")
		return
	}

	_, attErr := a.attrValueSrv.Create(c, &req)

	if attErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}
	if attErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}

	sessions.Set(c, "message", custom_messages.AttValueCreatedSuccessfully)
	c.Redirect(http.StatusFound, "/admins/attribute-values/create")
	return
}

func (a *AdminHandler) IndexAttributeValues(c *gin.Context) {

	attributes, err := a.attrValueSrv.IndexAttribute(c)
	if err.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}
	if err.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}

	html.Render(c, 200, "admin_index_attribute_value", gin.H{
		"TITLE":      "index attribute-value",
		"ATTRIBUTES": attributes,
	})

}

func (a *AdminHandler) ShowAttributeValues(c *gin.Context) {
	attributeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	attributeData, oldErr := a.attributeSrv.Show(c, attributeID)

	if oldErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}
	if oldErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}

	html.Render(c, 200, "admin_show_attribute_values_index", gin.H{
		"TITLE":     "Show values of attribute ",
		"ATTRIBUTE": attributeData,
	})
	return
}

func (a *AdminHandler) EditAttributeValues(c *gin.Context) {
	attValID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}

	attributes, attrsErr := a.attributeSrv.Index(c)
	if attrsErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}
	if attrsErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}

	attributeValueData, avErr := a.attrValueSrv.Show(c, attValID)

	if avErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}
	if avErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}

	html.Render(c, 200, "admin_show_attribute_values_edit", gin.H{
		"TITLE":          "edit attribute-values of a attribute ",
		"ATTRIBUTEVALUE": attributeValueData,
		"ATTRIBUTES":     attributes,
	})
	return

}

func (a *AdminHandler) UpdateAttributeValues(c *gin.Context) {

	//convert string id to int
	attValID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}

	url := fmt.Sprintf("/admins/attribute-values/%d/edit", attValID)

	//bind data
	var req requests.UpdateAttributeValueRequest
	_ = c.Request.ParseForm()
	bindErr := c.ShouldBind(&req)
	if bindErr != nil {
		errors.Init()
		errors.SetFromErrors(err)
		if req.AttributeID <= 0 {
			errors.Add("attribute_id", custom_messages.SelectOne)
		}
		fmt.Println("validation failed : ", err.Error())
		fmt.Println("attribute id : ", req.AttributeID)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, url)
		return
	}

	//find attribute-value
	_, oldErr := a.attrValueSrv.Show(c, attValID)
	if oldErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}
	if oldErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/attribute-values")
		return
	}

	//update attribute-value
	updateErr := a.attrValueSrv.Update(c, attValID, &req)
	if updateErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, url)
		return
	}
	if updateErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, url)
		return
	}

	//update was successful
	sessions.Set(c, "message", custom_messages.AttributeValueUpdateSuccessfully)
	c.Redirect(http.StatusFound, "/admins/attribute-values")
	return

}
