package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"shop/internal/modules/admin/requests"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/custom_messages"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/old"
	"shop/internal/pkg/sessions"
	"shop/internal/pkg/util"
	"slices"
	"strconv"
	"strings"
)

func (a *AdminHandler) IndexProduct(c *gin.Context) {
	products, err := a.productSrv.Index(c)
	if err.Code == 404 {
		c.JSON(http.StatusOK, gin.H{"msg": "محصولی یافت نشد / محصول اضافه کنید"})
		return
	}
	if err.Code == 500 {
		html.Error500(c)
		return
	}

	html.Render(c, http.StatusFound, "modules/admin/html/admin_index_product", gin.H{
		"TITLE":    "لیست محصولات",
		"PRODUCTS": products,
	})
	return
}

func (a *AdminHandler) CreateProduct(c *gin.Context) {
	categories, err := a.categorySrv.GetAllCategories(c)
	if err.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if err.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if len(categories.Data) == 0 {
		sessions.Set(c, "message", custom_messages.ThereIsNoAnyCategories)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	brands, _ := a.brandSrv.Index(c)
	if len(brands.Data) == 0 {
		sessions.Set(c, "message", custom_messages.ThereIsNoAnyBrand)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	html.Render(c, http.StatusFound, "modules/admin/html/admin_create_product", gin.H{
		"TITLE":      "ایجاد محصول | توجه: دسته بندی روت انتخاب نشود!",
		"CATEGORIES": categories,
		"BRANDS":     brands.Data,
	})
	return
}

func (a *AdminHandler) StoreProduct(c *gin.Context) {
	var req requests.CreateProductRequest
	_ = c.Request.ParseForm()
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

	//check uniqueness of sku
	IsUnique, CheckErr := a.productSrv.CheckSkuIsUnique(c, req.Sku)
	if CheckErr.Code == 500 {
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		sessions.Set(c, "message", custom_error.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}
	if !IsUnique {
		errors.Init()
		errors.Add("sku", custom_error.MustBeUnique)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}

	//check category_id existence
	_, cErr := a.categorySrv.Show(c, req.CategoryID)
	if cErr.Code > 0 {
		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		sessions.Set(c, "message", "شناسه کتگوری نامعتبر است")
		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}

	imagesForm, _ := c.MultipartForm()
	imagesFile := imagesForm.File["images[]"]
	//check required validation
	if imagesFile == nil {
		errors.Init()
		errors.Add("images", custom_error.IsRequired)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}

	var imagesStoredPath []string
	for _, image := range imagesFile {
		extension := filepath.Ext(image.Filename)

		// file extension validation
		ok := slices.Contains(util.AllowImageExtensions(), extension)
		if !ok {
			errors.Init()
			errors.Add("images", custom_error.MustBeImage)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, "/admins/products/create")
			return

		}

		//generate file name
		imageGenerateFileName := util.GenerateFilename(image.Filename)
		imagesStoredPath = append(imagesStoredPath, imageGenerateFileName)

		//check upload and store file to storage bucket
		if os.Getenv("STORAGE_STATUS") == "active" {
			go func() {
				{
					imageFile, _ := image.Open()
					err := a.dep.Storage.UploadFile(imageFile, os.Getenv("STORAGE_PRODUCT_PATH")+imageGenerateFileName)
					if err != nil {
						log.Println("-- failed to upload file to s3 : ", err)
					}
				}
			}()
		}

		//store images on disk
		saveUploadedImage := c.SaveUploadedFile(image, viper.GetString("Upload.Products")+imageGenerateFileName)
		if saveUploadedImage != nil {
			for _, imageStorePath := range imagesStoredPath {
				_ = os.Remove(viper.GetString("Upload.Products") + imageStorePath)
			}

			errors.Init()
			errors.Add("images", custom_error.StoreImageOnDiskFailed)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, "/admins/products/create")
			return
		}
	}
	req.ProductImage = imagesStoredPath
	_, pErr := a.productSrv.Create(c, &req)
	if pErr.Code > 0 {
		//remove images from disk
		for _, img := range imagesStoredPath {
			_ = os.Remove(viper.GetString("Upload.Categories") + img)
		}

		sessions.Set(c, "message", custom_error.SomethingWrongHappened)
		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}

	sessions.Set(c, "message", custom_error.SuccessfullyCreated)
	c.Redirect(http.StatusFound, "/admins/products/create")
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
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if pErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	html.Render(c, http.StatusFound, "admin_show_product",
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
	productShow, allProductsInMongo, pErr := a.productSrv.Show(c, "id", pID)

	if pErr.Code == 404 {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if pErr.Code == 500 {
		html.Error500(c)
		return
	}

	categories, cErr := a.categorySrv.GetAllCategories(c)
	if cErr.Code == 404 {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if cErr.Code == 500 {
		html.Error500(c)
		return
	}

	brands, bErr := a.brandSrv.Index(c)
	if bErr.Code == 404 {
		sessions.Set(c, "message", "برند یافت نشد")
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if bErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	recommendations, _ := a.productSrv.FetchAllRecommendation(c, pID)

	html.Render(c, http.StatusFound, "modules/admin/html/admin_edit_product",
		gin.H{
			"TITLE":           "ویرایش محصول",
			"PRODUCT":         productShow,
			"AllProducts":     allProductsInMongo, // products in collection products in mongodb
			"CATEGORIES":      categories,
			"BRANDS":          brands,
			"RECOMMENDATIONS": recommendations,
		},
	)
	return
}

func (a *AdminHandler) UpdateProduct(c *gin.Context) {

	//convert string id to int
	productID, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	//var req
	var req requests.UpdateProductRequest
	_ = c.Request.ParseForm()

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

	//select product from db
	selectedProduct, _, pErr := a.productSrv.Show(c, "id", productID)
	if pErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if pErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	//check brand id existence
	if selectedProduct.BrandID != req.BrandID {
		//selectedBrand
		_, bErr := a.brandSrv.Show(c, int(req.BrandID))
		if bErr.Code == 404 {
			sessions.Set(c, "message", "شناسه برند نامعتبر می باشد")
			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}
		if bErr.Code == 500 {
			sessions.Set(c, "message", custom_error.InternalServerError)
			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}
	}

	//check category id existence
	if selectedProduct.CategoryID != uint(req.CategoryID) {
		//selectedCategory
		_, cErr := a.categorySrv.Show(c, req.CategoryID)
		if cErr.Code == 404 {
			sessions.Set(c, "message", "شناسه کتگوری نامعتبر می باشد")
			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}
		if cErr.Code == 500 {
			sessions.Set(c, "message", custom_error.InternalServerError)
			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}
	}

	//check uniqueness of sku
	if selectedProduct.Sku != strings.TrimSpace(req.Sku) {
		IsUnique, CheckErr := a.productSrv.CheckSkuIsUnique(c, req.Sku)
		if CheckErr.Code == 500 {
			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			sessions.Set(c, "message", custom_error.SomethingWrongHappened)
			c.Redirect(http.StatusFound, url)
			return
		}
		if !IsUnique {
			fmt.Println("unique error :", CheckErr)
			errors.Init()
			errors.Add("sku", custom_error.MustBeUnique)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}
	}

	//update product
	updateErr := a.productSrv.Update(c, productID, &req)
	if updateErr.Code > 0 {
		sessions.Set(c, "message", custom_messages.ProductUpdateFailed)
		c.Redirect(http.StatusFound, "/admins/products/")
	}

	//--add product recommendations

	// get recommendation IDs from Form
	productRecommendationIDs := c.PostFormArray("recommendations")

	go a.productSrv.AddRecommendation(c, productID, productRecommendationIDs)

	sessions.Set(c, "message", custom_messages.ProductUpdatedSuccessfully)
	c.Redirect(http.StatusFound, url)
	return

}
