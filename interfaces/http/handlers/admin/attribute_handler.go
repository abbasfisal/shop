package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/domain/domain_err"
	"shop/infrastructure/messages"
	"shop/interfaces/http/requests/admin"
	"shop/interfaces/http/response"
	"shop/pkg/errors"
	"shop/pkg/old"
	"shop/pkg/sessions"
	"strconv"
	"strings"
)

func (a *AdminHandler) CreateAttribute(c *gin.Context) {
	response.Render(c, http.StatusFound, "admin_create_attribute",
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
		response.Error500(c)
		return
	}
	response.Render(c, 200, "admin_index_attribute",
		gin.H{
			"TITLE":      "لیست اتریبیوت ها",
			"ATTRIBUTES": attributes,
		})
	return

}

func (a *AdminHandler) ShowAttribute(c *gin.Context) {
	attributeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}

	attributeShow, attErr := a.attributeSrv.Show(c, attributeID)
	if attErr.Code == 404 {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}
	if attErr.Code == 500 {
		response.Error500(c)
		return
	}

	response.Render(c, http.StatusFound, "admin_show_attribute",
		gin.H{
			"TITLE":     "نمایش اتریبیوت",
			"ATTRIBUTE": attributeShow,
		},
	)
}

func (a *AdminHandler) UpdateAttribute(c *gin.Context) {

	attributeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}

	url := fmt.Sprintf("/admins/attributes/%d/edit", attributeID)

	oldAttribute, oldErr := a.attributeSrv.Show(c, attributeID)

	if oldErr.Code == 404 {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}
	if oldErr.Code == 500 {
		sessions.Set(c, "message", domain_err.InternalServerError)
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
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}
	if updateErr.Code == 500 {
		sessions.Set(c, "message", domain_err.InternalServerError)
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
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
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
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}
	if pAttrErr.Code == 500 {
		sessions.Set(c, "message", domain_err.InternalServerError)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	sessions.Set(c, "message", domain_err.SuccessfullyCreated)

	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}

func (a *AdminHandler) GetAttributesByCategoryID(c *gin.Context) {
	// legacy AJAX route used by the attribute-value pages to fill the picker.
	// Attributes are not category scoped anymore → return every attribute
	// (with its values) in the shape those pages expect: {Data:[{ID,Title}]}.
	attributes, aErr := a.attributeSrv.Index(c)
	if aErr.Code > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attributes"})
		return
	}
	c.JSON(http.StatusOK, attributes)
}

// GetAttributesJSON feeds the product create/edit combination builder:
// GET /admins/api/attributes → {data:[{id,title,code,values:[{id,value}]}]}
func (a *AdminHandler) GetAttributesJSON(c *gin.Context) {
	attributes, aErr := a.attributeSrv.Index(c)
	if aErr.Code > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain_err.SomethingWrongHappened})
		return
	}

	type valueJSON struct {
		ID    uint   `json:"id"`
		Value string `json:"value"`
	}
	type attributeJSON struct {
		ID     uint        `json:"id"`
		Title  string      `json:"title"`
		Code   string      `json:"code"`
		Values []valueJSON `json:"values"`
	}

	data := make([]attributeJSON, 0, len(attributes.Data))
	for _, attr := range attributes.Data {
		row := attributeJSON{ID: attr.ID, Title: attr.Title, Code: attr.Code, Values: []valueJSON{}}
		if attr.AttributeValues != nil {
			for _, v := range attr.AttributeValues.Data {
				row.Values = append(row.Values, valueJSON{ID: v.ID, Value: v.Title})
			}
		}
		data = append(data, row)
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
