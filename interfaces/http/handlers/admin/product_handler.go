package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

// flashFormErrors stores Persian validation errors + old input so the
// create/edit page can re-render the form exactly as it was submitted.
func flashFormErrors(c *gin.Context, formErrs map[string]string) {
	errors.Init()
	for key, msg := range formErrs {
		errors.Add(key, msg)
	}
	sessions.Set(c, "errors", errors.ToString())

	old.Init()
	old.Set(c)
	sessions.Set(c, "olds", old.ToString())
}

// oldVariantsJSON rebuilds the posted variants[i][...] rows as JSON so a
// failed validation does not wipe the combination table (read from the
// non-destructive session copy, WithGlobalData flashes the rest).
func oldVariantsJSON(c *gin.Context) string {
	raw := sessions.GET(c, "olds")
	if raw == "" {
		return "[]"
	}
	var form map[string][]string
	if err := json.Unmarshal([]byte(raw), &form); err != nil {
		return "[]"
	}
	rows := requests.ParseVariantRows(form)
	out, err := json.Marshal(rows)
	if err != nil {
		return "[]"
	}
	return string(out)
}

// oldFields flattens the flashed old input to one value per key — used by the
// create/edit templates for <select> / <radio> state after a failed post.
func oldFields(c *gin.Context) map[string]string {
	raw := sessions.GET(c, "olds")
	out := map[string]string{}
	if raw == "" {
		return out
	}
	var form map[string][]string
	if err := json.Unmarshal([]byte(raw), &form); err != nil {
		return out
	}
	for key, vals := range form {
		if len(vals) > 0 {
			out[key] = vals[0]
		}
	}
	return out
}

// productListQuery reads the admin list filters (q, status, category_id,
// in_stock, attr[<code>][] and sort) out of the query string.
func productListQuery(c *gin.Context) requests.ProductListQuery {
	q := requests.ProductListQuery{
		Q:          strings.TrimSpace(c.Query("q")),
		Status:     strings.TrimSpace(c.Query("status")),
		CategoryID: util.StringToUint(c.Query("category_id")),
		InStock:    c.Query("in_stock") == "1",
		Sort:       c.Query("sort"),
	}

	attrs := map[string][]uint{}
	for key, values := range c.Request.URL.Query() {
		if !strings.HasPrefix(key, "attr[") {
			continue
		}
		// attr[<code>] or attr[<code>][] → <code>
		code := strings.TrimPrefix(key, "attr[")
		code = strings.TrimSuffix(code, "[]")
		code = strings.TrimSuffix(code, "]")
		if code == "" {
			continue
		}
		ids := make([]uint, 0, len(values))
		for _, v := range values {
			if id := util.StringToUint(v); id > 0 {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			attrs[code] = ids
		}
	}
	q.Attrs = attrs

	return q
}

// chosenAttrs turns the parsed filter into {code: {valueID: true}} so the
// template can tick the active attribute checkboxes with index/map lookups.
func chosenAttrs(attrs map[string][]uint) map[string]map[uint]bool {
	out := map[string]map[uint]bool{}
	for code, ids := range attrs {
		set := map[uint]bool{}
		for _, id := range ids {
			set[id] = true
		}
		out[code] = set
	}
	return out
}

func (a *AdminHandler) IndexProduct(c *gin.Context) {
	query := productListQuery(c)

	products, err := a.productSrv.Index(c, query)
	if err.Code == 404 {
		products = &responses.Products{}
	} else if err.Code == 500 {
		response.Error500(c)
		return
	}

	// filter sidebar data (Laravel index.blade.php)
	categories, _ := a.categorySrv.GetAllCategories(c)
	if categories == nil {
		categories = &responses.Categories{}
	}
	attributes, _ := a.attributeSrv.Index(c)
	if attributes == nil {
		attributes = &responses.Attributes{}
	}

	response.Render(c, http.StatusOK, "modules/admin/html/admin_index_product", gin.H{
		"TITLE":      "لیست محصولات",
		"PRODUCTS":   products,
		"CATEGORIES": categories,
		"ATTRIBUTES": attributes,
		"FILTER": gin.H{
			"Q":          query.Q,
			"Status":     query.Status,
			"CategoryID": query.CategoryID,
			"InStock":    query.InStock,
			"Sort":       query.Sort,
		},
		"CHOSEN": chosenAttrs(query.Attrs),
	})
	return
}

func (a *AdminHandler) CreateProduct(c *gin.Context) {
	categories, err := a.categorySrv.GetAllCategories(c)
	if err.Code == 404 {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if err.Code == 500 {
		sessions.Set(c, "message", domain_err.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if len(categories.Data) == 0 {
		sessions.Set(c, "message", custom_messages.ThereIsNoAnyCategories)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	brands, _ := a.brandSrv.Index(c)
	if brands == nil || len(brands.Data) == 0 {
		sessions.Set(c, "message", custom_messages.ThereIsNoAnyBrand)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	response.Render(c, http.StatusOK, "modules/admin/html/admin_create_product", gin.H{
		"TITLE":       "ایجاد محصول | توجه: دسته بندی روت انتخاب نشود!",
		"CATEGORIES":  categories,
		"BRANDS":      brands.Data,
		"OLD":         oldFields(c),
		"OLDVARIANTS": oldVariantsJSON(c),
	})
	return
}

func (a *AdminHandler) StoreProduct(c *gin.Context) {
	_ = c.Request.ParseMultipartForm(32 << 20)

	var req requests.CreateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}

	// Laravel StoreProductRequest rules (type/status/variants/simple fields)
	if formErrs := requests.ValidateProductForm(c.Request.PostForm); len(formErrs) > 0 {
		flashFormErrors(c, formErrs)
		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}
	req.Variants = requests.ParseVariantRows(c.Request.PostForm)

	//check uniqueness of sku
	IsUnique, CheckErr := a.productSrv.CheckSkuIsUnique(c, req.Sku)
	if CheckErr.Code == 500 {
		flashFormErrors(c, map[string]string{"sku": domain_err.SomethingWrongHappened})
		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}
	if !IsUnique {
		flashFormErrors(c, map[string]string{"sku": domain_err.MustBeUnique})
		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}

	//check category_id existence
	_, cErr := a.categorySrv.Show(c, req.CategoryID)
	if cErr.Code > 0 {
		flashFormErrors(c, map[string]string{"category_id": "شناسه کتگوری نامعتبر است"})
		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}

	imagesForm, _ := c.MultipartForm()
	var imagesFile []*multipart.FileHeader
	if imagesForm != nil {
		imagesFile = imagesForm.File["images[]"]
	}

	var imagesStoredPath []string
	for _, image := range imagesFile {
		extension := filepath.Ext(image.Filename)

		// file extension validation
		ok := slices.Contains(util.AllowImageExtensions(), extension)
		if !ok {
			flashFormErrors(c, map[string]string{"images": domain_err.MustBeImage})
			c.Redirect(http.StatusFound, "/admins/products/create")
			return
		}

		//generate file name
		imageGenerateFileName := util.GenerateFilename(image.Filename)
		imagesStoredPath = append(imagesStoredPath, imageGenerateFileName)

		//check upload and store file to storage bucket
		if os.Getenv("STORAGE_STATUS") == "active" {
			go func(fileName string, fileHeader *multipart.FileHeader) {
				imageFile, _ := fileHeader.Open()
				if err := a.dep.Storage.UploadFile(imageFile, os.Getenv("STORAGE_PRODUCT_PATH")+fileName); err != nil {
					log.Println("-- failed to upload file to s3 : ", err)
				}
			}(imageGenerateFileName, image)
		}

		//store images on disk
		saveUploadedImage := c.SaveUploadedFile(image, viper.GetString("Upload.Products")+imageGenerateFileName)
		if saveUploadedImage != nil {
			for _, imageStorePath := range imagesStoredPath {
				_ = os.Remove(viper.GetString("Upload.Products") + imageStorePath)
			}

			flashFormErrors(c, map[string]string{"images": domain_err.StoreImageOnDiskFailed})
			c.Redirect(http.StatusFound, "/admins/products/create")
			return
		}
	}
	req.ProductImage = imagesStoredPath

	_, pErr := a.productSrv.Create(c, &req)
	if pErr.Code > 0 {
		//remove images from disk
		for _, img := range imagesStoredPath {
			_ = os.Remove(viper.GetString("Upload.Products") + img)
		}

		sessions.Set(c, "message", domain_err.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}

	sessions.Set(c, "message", custom_messages.ProductCreatedSuccessfully)
	c.Redirect(http.StatusFound, "/admins/products")
	return
}

func (a *AdminHandler) ShowProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	selectedP, _, pErr := a.productSrv.Show(c, "id", productID)
	if pErr.Code == 404 {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if pErr.Code == 500 {
		sessions.Set(c, "message", domain_err.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	response.Render(c, http.StatusOK, "admin_show_product",
		gin.H{
			"TITLE":   "نمایش محصول",
			"PRODUCT": selectedP,
		},
	)
}

func (a *AdminHandler) EditProduct(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	// get product data by repo
	productShow, allProductBriefs, pErr := a.productSrv.Show(c, "id", pID)

	if pErr.Code == 404 {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if pErr.Code == 500 {
		response.Error500(c)
		return
	}

	categories, cErr := a.categorySrv.GetAllCategories(c)
	if cErr.Code == 404 {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if cErr.Code == 500 {
		response.Error500(c)
		return
	}

	brands, bErr := a.brandSrv.Index(c)
	if bErr.Code == 404 {
		sessions.Set(c, "message", "برند یافت نشد")
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if bErr.Code == 500 {
		sessions.Set(c, "message", domain_err.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	recommendations, _ := a.productSrv.FetchAllRecommendation(c, pID)

	// first variant (zero value when the product has none) — keeps the simple
	// form section free of out-of-range index in the template
	firstVariant := responses.ProductInventory{}
	if productShow != nil && productShow.ProductInventories != nil && len(productShow.ProductInventories.Data) > 0 {
		firstVariant = productShow.ProductInventories.Data[0]
	}

	response.Render(c, http.StatusOK, "modules/admin/html/admin_edit_product",
		gin.H{
			"TITLE":           "ویرایش محصول",
			"PRODUCT":         productShow,
			"AllProducts":     allProductBriefs, // lightweight product summaries for recommendation picker
			"CATEGORIES":      categories,
			"BRANDS":          brands.Data,
			"RECOMMENDATIONS": recommendations,
			"OLD":             oldFields(c),
			"OLDVARIANTS":     oldVariantsJSON(c),
			"INV":             firstVariant,
		},
	)
	return
}

func (a *AdminHandler) UpdateProduct(c *gin.Context) {

	//convert string id to int
	productID, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		sessions.Set(c, "message", domain_err.IDIsNotCorrect)

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	//var req
	var req requests.UpdateProductRequest
	_ = c.Request.ParseMultipartForm(32 << 20)

	//var url
	url := fmt.Sprintf("/admins/products/%d/edit", productID)

	//bind request
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, url)
		return
	}

	// Laravel StoreProductRequest rules (type/status/variants/simple fields)
	if formErrs := requests.ValidateProductForm(c.Request.PostForm); len(formErrs) > 0 {
		flashFormErrors(c, formErrs)
		c.Redirect(http.StatusFound, url)
		return
	}
	req.Variants = requests.ParseVariantRows(c.Request.PostForm)

	//select product from db
	selectedProduct, _, pErr := a.productSrv.Show(c, "id", productID)
	if pErr.Code == 404 {
		sessions.Set(c, "message", domain_err.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if pErr.Code == 500 {
		sessions.Set(c, "message", domain_err.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	//check brand id existence
	if selectedProduct.BrandID != req.BrandID {
		//selectedBrand
		_, bErr := a.brandSrv.Show(c, int(req.BrandID))
		if bErr.Code == 404 {
			flashFormErrors(c, map[string]string{"brand_id": "شناسه برند نامعتبر می باشد"})
			c.Redirect(http.StatusFound, url)
			return
		}
		if bErr.Code == 500 {
			sessions.Set(c, "message", domain_err.InternalServerError)
			c.Redirect(http.StatusFound, url)
			return
		}
	}

	//check category id existence
	if selectedProduct.CategoryID != uint(req.CategoryID) {
		//selectedCategory
		_, cErr := a.categorySrv.Show(c, req.CategoryID)
		if cErr.Code == 404 {
			flashFormErrors(c, map[string]string{"category_id": "شناسه کتگوری نامعتبر می باشد"})
			c.Redirect(http.StatusFound, url)
			return
		}
		if cErr.Code == 500 {
			sessions.Set(c, "message", domain_err.InternalServerError)
			c.Redirect(http.StatusFound, url)
			return
		}
	}

	//check uniqueness of sku
	if selectedProduct.Sku != strings.TrimSpace(req.Sku) {
		IsUnique, CheckErr := a.productSrv.CheckSkuIsUnique(c, req.Sku)
		if CheckErr.Code == 500 {
			sessions.Set(c, "message", domain_err.SomethingWrongHappened)
			c.Redirect(http.StatusFound, url)
			return
		}
		if !IsUnique {
			flashFormErrors(c, map[string]string{"sku": domain_err.MustBeUnique})
			c.Redirect(http.StatusFound, url)
			return
		}
	}

	//update product (+ its variant rows)
	updateErr := a.productSrv.Update(c, productID, &req)
	if updateErr.Code > 0 {
		sessions.Set(c, "message", custom_messages.ProductUpdateFailed)
		c.Redirect(http.StatusFound, url)
		return
	}

	//--add product recommendations
	// get recommendation IDs from Form
	productRecommendationIDs := c.PostFormArray("recommendations")

	go a.productSrv.AddRecommendation(c, productID, productRecommendationIDs)

	sessions.Set(c, "message", custom_messages.ProductUpdatedSuccessfully)
	c.Redirect(http.StatusFound, url)
	return
}
