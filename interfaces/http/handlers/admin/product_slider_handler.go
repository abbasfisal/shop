package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	responses "shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/infrastructure/messages"
	"shop/interfaces/http/requests/admin"
	"shop/interfaces/http/response"
	"shop/pkg/errors"
	"shop/pkg/old"
	"shop/pkg/sessions"
	"shop/pkg/util"
)

// sliderFormData binds + validates the slider form (create and edit).
// It returns false when the request must stop (errors are flashed).
func (a *AdminHandler) sliderFormData(c *gin.Context) (requests.CreateProductSliderRequest, bool) {
	_ = c.Request.ParseForm()

	var req requests.CreateProductSliderRequest
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		return req, false
	}

	if formErrs := requests.ValidateProductSliderForm(c.Request.PostForm); len(formErrs) > 0 {
		errors.Init()
		for key, msg := range formErrs {
			errors.Add(key, msg)
		}
		sessions.Set(c, "errors", errors.ToString())
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())
		return req, false
	}

	req.ProductIDs = requests.ParseSliderProductIDs(c.Request.PostForm)
	req.Position = strings.TrimSpace(req.Position)
	req.Status = requests.SliderStatus(c.Request.PostForm)
	req.Title = strings.TrimSpace(req.Title)
	return req, true
}

// sliderViewData is the shared payload of the create/edit slider pages.
func (a *AdminHandler) sliderViewData(c *gin.Context) gin.H {
	categories, _ := a.categorySrv.GetAllCategories(c)
	if categories == nil {
		categories = &responses.Categories{}
	}
	return gin.H{
		"POSITIONS":  responses.SliderPositions(),
		"CATEGORIES": categories,
	}
}

func (a *AdminHandler) IndexProductSlider(c *gin.Context) {
	sliders, err := a.slidersSrv.Index(c)
	if err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/sliders")
		return
	}

	response.Render(c, http.StatusOK, "admin_index_slider", gin.H{
		"TITLE":   "اسلایدرهای محصولات",
		"SLIDERS": responses.ToSliders(sliders),
	})
	return
}

func (a *AdminHandler) CreateProductSlider(c *gin.Context) {
	data := a.sliderViewData(c)
	data["TITLE"] = "ایجاد اسلایدر محصولات"
	data["SLIDER"] = responses.ProductSlider{Position: "after_categories", Status: "published"}
	data["SELECTED_JSON"] = "[]"

	response.Render(c, http.StatusOK, "admin_create_slider", data)
	return
}

func (a *AdminHandler) StoreProductSlider(c *gin.Context) {
	req, ok := a.sliderFormData(c)
	if !ok {
		c.Redirect(http.StatusFound, "/admins/sliders/create")
		return
	}

	slider, err := a.slidersSrv.Store(c, &req)
	if err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/sliders/create")
		return
	}

	sessions.Set(c, "message", "اسلایدر «"+slider.Title+"» ایجاد شد")
	c.Redirect(http.StatusFound, "/admins/sliders")
	return
}

func (a *AdminHandler) EditProductSlider(c *gin.Context) {
	sliderID := util.StringToUint(c.Param("id"))
	slider, err := a.slidersSrv.Show(c, sliderID)
	if err != nil {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/sliders")
		return
	}

	sliderView := responses.ToSlider(slider)
	selectedJSON, _ := json.Marshal(sliderView.Products)
	if selectedJSON == nil {
		selectedJSON = []byte("[]")
	}

	data := a.sliderViewData(c)
	data["TITLE"] = "ویرایش اسلایدر محصولات"
	data["SLIDER"] = sliderView
	data["SELECTED_JSON"] = string(selectedJSON)
	data["IS_EDIT"] = true

	response.Render(c, http.StatusOK, "admin_create_slider", data)
	return
}

func (a *AdminHandler) UpdateProductSlider(c *gin.Context) {
	sliderID := util.StringToUint(c.Param("id"))
	if sliderID == 0 {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/sliders")
		return
	}

	req, ok := a.sliderFormData(c)
	if !ok {
		c.Redirect(http.StatusFound, "/admins/sliders/"+c.Param("id")+"/edit")
		return
	}

	if err := a.slidersSrv.Update(c, sliderID, &req); err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/sliders/"+c.Param("id")+"/edit")
		return
	}

	sessions.Set(c, "message", custom_messages.ProductUpdatedSuccessfully)
	c.Redirect(http.StatusFound, "/admins/sliders")
	return
}

func (a *AdminHandler) DeleteProductSlider(c *gin.Context) {
	sliderID := util.StringToUint(c.Param("id"))
	if sliderID == 0 {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/sliders")
		return
	}

	if err := a.slidersSrv.Delete(c, sliderID); err != nil {
		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
	} else {
		sessions.Set(c, "message", custom_messages.DeleteSuccessfully)
	}
	c.Redirect(http.StatusFound, "/admins/sliders")
	return
}

// SearchProductsJSON is the AJAX product picker used by the slider form
// (GET /admins/api/products?q=...).
func (a *AdminHandler) SearchProductsJSON(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}

	products, err := a.productSrv.Index(c, requests.ProductListQuery{
		Q:    query,
		Sort: "name",
	})
	if err.Code > 0 || products == nil {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
		return
	}

	limit := 10
	if v := c.Query("limit"); v != "" {
		if n := util.StringToUint(v); n > 0 && n <= 50 {
			limit = int(n)
		}
	}

	type productJSON struct {
		ID         uint   `json:"id"`
		Title      string `json:"title"`
		Sku        string `json:"sku"`
		Slug       string `json:"slug"`
		Image      string `json:"image"`
		Price      uint   `json:"price"`
		InStock    bool   `json:"in_stock"`
		Status     string `json:"status"`
		StatusText string `json:"status_text"`
	}

	data := make([]productJSON, 0, limit)
	for _, p := range products.Data {
		if len(data) >= limit {
			break
		}
		image := ""
		if p.Images != nil && len(p.Images.Data) > 0 {
			image = p.Images.Data[0].OriginalPath
		}
		price := p.MinPrice
		if price == 0 {
			price = p.SalePrice
		}
		data = append(data, productJSON{
			ID:         p.ID,
			Title:      p.Title,
			Sku:        p.Sku,
			Slug:       p.Slug,
			Image:      image,
			Price:      price,
			InStock:    p.InStock,
			Status:     p.Status,
			StatusText: p.StatusText,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
	return
}
