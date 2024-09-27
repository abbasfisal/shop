package handlers

import "C"
import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/services/attribute"
	attributeValue "shop/internal/modules/admin/services/attribute_value"
	"shop/internal/modules/admin/services/auth"
	"shop/internal/modules/admin/services/brand"
	"shop/internal/modules/admin/services/category"
	"shop/internal/modules/admin/services/customer"
	"shop/internal/modules/admin/services/product"
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

type AdminHandler struct {
	authSrv      auth.AuthenticateServiceInterface
	categorySrv  category.CategoryServiceInterface
	productSrv   product.ProductServiceInterface
	attributeSrv attribute.AttributeServiceInterface
	attrValueSrv attributeValue.AttributeValueServiceInterface
	brandSrv     brand.BrandServiceInterface
	customerSrv  customer.CustomerServiceInterface

	i18nBundle *i18n.Bundle
	//order service
	//user service
	//cart service
}

func NewAdminHandler(
	authSrv auth.AuthenticateServiceInterface,
	categorySrv category.CategoryService,
	productSrv product.ProductServiceInterface,
	attributeSrv attribute.AttributeServiceInterface,
	attrValueSrv attributeValue.AttributeValueServiceInterface,
	brandSrv brand.BrandServiceInterface,
	customerSrv customer.CustomerServiceInterface,

	i18nBundle *i18n.Bundle,
) AdminHandler {
	return AdminHandler{
		authSrv:      authSrv,
		categorySrv:  categorySrv,
		productSrv:   productSrv,
		attributeSrv: attributeSrv,
		attrValueSrv: attrValueSrv,
		brandSrv:     brandSrv,
		customerSrv:  customerSrv,

		i18nBundle: i18nBundle,
	}
}

func (a AdminHandler) ShowLogin(c *gin.Context) {
	html.Render(c, http.StatusOK, "modules/admin/html/admin_login", gin.H{"title": "login"})
	return
}

func (a AdminHandler) PostLogin(c *gin.Context) {
	var req requests.LoginRequest
	c.Request.ParseForm()

	if err := c.ShouldBind(&req); err != nil {
		errors.SetErrors(c, a.i18nBundle, err)

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/login")
		return
	}

	user, loginErr := a.authSrv.Login(context.TODO(), req)
	if loginErr.Error() != "" {
		if loginErr.Code == 404 {
			html.Render(c, http.StatusFound, "modules/admin/html/admin_login", gin.H{
				"MESSAGE": loginErr.Error(),
			})
			return
		}
		if loginErr.Code == 500 {
			html.Error500(c)
			return
		}
	}

	sessions.Set(c, "auth_id", strconv.Itoa(int(user.ID)))
	c.Redirect(http.StatusFound, "/admins/home")
}

func (a AdminHandler) ShowHome(c *gin.Context) {
	html.Render(c, http.StatusOK, "modules/admin/html/admin_home",
		gin.H{
			"TITLE": "admin home page",
		})
	return
}

//-------------------------------
//		category routes
//-------------------------------

func (a AdminHandler) IndexCategory(c *gin.Context) {

	categories, err := a.categorySrv.Index(c)

	if err.Code == 400 {
		c.JSON(200, gin.H{
			"data": "empty",
		})
	} else if err.Code == 500 {
		html.Error500(c)
		return
	}
	html.Render(c, 200, "modules/admin/html/admin_index_category", gin.H{
		"TITLE":      "Index Category",
		"CATEGORIES": categories,
	})
	return
}

func (a AdminHandler) CreateCategory(c *gin.Context) {
	categories, _ := a.categorySrv.GetAllCategories(c)
	html.Render(c, http.StatusFound, "modules/admin/html/admin_create_category", gin.H{
		"TITLE":      "create new category",
		"CATEGORIES": categories,
	})
	return
}

func (a AdminHandler) StoreCategory(c *gin.Context) {
	var req requests.CreateCategoryRequest
	_ = c.Request.ParseForm()
	if err := c.ShouldBind(&req); err != nil {
		errors.Init()
		errors.SetFromErrors(err)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/categories/create")
		return
	}

	//slug unique validation
	if ok := a.categorySrv.CheckSlugUniqueness(context.TODO(), req.Slug); ok {
		errors.Init()
		errors.Add("slug", custom_error.MustBeUnique)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/categories/create")
		return
	}

	//------------------------
	//	upload and save image
	//------------------------
	imageFile, imageErr := c.FormFile("image")
	//require validation on image
	//if imageErr != nil {
	//	//validation on required tag
	//	errors.Init()
	//	errors.Add("image", custom_error.IsRequired)
	//	sessions.Set(c, "errors", errors.ToString())
	//
	//	old.Init()
	//	old.Set(c)
	//	sessions.Set(c, "olds", old.ToString())
	//
	//	c.Redirect(http.StatusFound, "/admins/categories/create")
	//
	//	return
	//}

	pathToUpload := ""
	if imageErr == nil {
		// file extension validation
		fileExtension := filepath.Ext(imageFile.Filename)
		ok := slices.Contains(util.AllowImageExtensions(), fileExtension)
		if !ok {
			errors.Init()
			errors.Add("image", custom_error.MustBeImage)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, "/admins/categories/create")
			return
		}

		//upload image and store on disk
		newImageName := util.GenerateFilename(imageFile.Filename)
		pathToUpload := viper.GetString("Upload.Categories") + newImageName
		uploadErr := c.SaveUploadedFile(imageFile, pathToUpload)
		if uploadErr != nil {
			fmt.Println("upload error:", uploadErr)
			errors.Init()
			errors.Add("image", custom_error.UploadImageError)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, "/admins/categories/create")

			return
		}
		req.Image = newImageName
	}

	//---------------------------
	//	end upload and save image
	//---------------------------

	newCategory, err := a.categorySrv.Create(context.TODO(), req)
	if err != nil || newCategory.ID <= 0 {
		_ = os.Remove(pathToUpload)
		fmt.Println("error in creating category : ", err)
		sessions.Set(c, "message", "خطا در ایجاد دسته بندی")
		c.Redirect(http.StatusFound, "/admins/categories/create")
		return
	}

	sessions.Set(c, "message", "ایجاد دسته بندی با موفقیت انجام شد")
	c.Redirect(http.StatusFound, "/admins/categories/create")
	return
}

func (a AdminHandler) ShowCategory(c *gin.Context) {

	catID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		catID = 1
	}
	cat, catErr := a.categorySrv.Show(context.TODO(), catID)
	if catErr.Code == 404 {
		html.Render(c, http.StatusOK, "modules/admin/html/admin_index_category", gin.H{
			"MESSAGE": custom_error.RecordNotFound,
		})
		return
	}
	if catErr.Code == 500 {
		html.Error500(c)
		return
	}

	html.Render(c, http.StatusFound, "modules/admin/html/admin_show_category", gin.H{
		"TITLE":    "show a category",
		"CATEGORY": cat,
	})
	return
}

func (a AdminHandler) EditCategory(c *gin.Context) {
	catID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}

	cat, cErr := a.categorySrv.Show(context.TODO(), catID)
	categories, _ := a.categorySrv.GetAllCategories(c)

	if cErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}
	if cErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}

	html.Render(c, http.StatusFound, "admin_edit_category", gin.H{
		"TITLE":      "edit category",
		"CATEGORY":   cat,
		"CATEGORIES": categories,
	})
	return
}

func (a AdminHandler) UpdateCategory(c *gin.Context) {

	catID, catIDErr := strconv.Atoi(c.Param("id"))
	if catIDErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}

	url := fmt.Sprintf("/admins/categories/%d/edit", catID)

	var req requests.UpdateCategoryRequest
	_ = c.Request.ParseForm()
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("------- bind err : ", err)
		errors.Init()
		errors.SetFromErrors(err)

		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, url)
		return
	}

	categoryData, catErr := a.categorySrv.Show(c, catID)
	if catErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}
	if catErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}

	//slug unique validation
	if categoryData.Slug != req.Slug {
		if ok := a.categorySrv.CheckSlugUniqueness(context.TODO(), req.Slug); ok {
			errors.Init()
			errors.Add("slug", custom_error.MustBeUnique)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}
	}

	//------------------------
	//	upload and save image
	//------------------------
	imageFile, imageErr := c.FormFile("image")

	pathToUpload := ""
	if imageErr != nil {
		req.Image = categoryData.Image
	}
	if imageErr == nil {
		// file extension validation
		fileExtension := filepath.Ext(imageFile.Filename)
		ok := slices.Contains(util.AllowImageExtensions(), fileExtension)
		if !ok {
			errors.Init()
			errors.Add("image", custom_error.MustBeImage)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)
			return
		}

		//upload image and store on disk
		newImageName := util.GenerateFilename(imageFile.Filename)
		pathToUpload = viper.GetString("Upload.Categories") + newImageName
		uploadErr := c.SaveUploadedFile(imageFile, pathToUpload)
		if uploadErr != nil {
			fmt.Println("upload error:", uploadErr)
			errors.Init()
			errors.Add("image", custom_error.UploadImageError)
			sessions.Set(c, "errors", errors.ToString())

			old.Init()
			old.Set(c)
			sessions.Set(c, "olds", old.ToString())

			c.Redirect(http.StatusFound, url)

			return
		}
		req.Image = newImageName
	}

	//---------------------------
	//	end upload and save image
	//---------------------------

	categoryUpdateErr := a.categorySrv.Edit(c, catID, req)
	fmt.Println("----- hanlder category edit : error  ", categoryUpdateErr)
	if categoryUpdateErr.Code > 0 {
		fmt.Println("------- update category err : ", categoryUpdateErr.Error())
		if categoryUpdateErr.Code == 404 {
			sessions.Set(c, "message", custom_error.RecordNotFound)
			c.Redirect(http.StatusFound, "/admins/categories/")
			return
		}
		if categoryUpdateErr.Code == 500 {
			sessions.Set(c, "message", custom_error.InternalServerError)
			c.Redirect(http.StatusFound, "/admins/categories/")
			return
		}
	}

	if pathToUpload != "" {
		delImg := os.Remove(viper.GetString("upload.categories") + categoryData.Image)
		if delImg != nil {
			fmt.Println("----- dele image  err: ", delImg.Error())
		}
		fmt.Println("--- delete succ imag")
	}
	sessions.Set(c, "message", custom_messages.CategoryUpdatedSucc)
	c.Redirect(http.StatusFound, "/admins/categories/")
	return

}

func (a AdminHandler) CategoryProducts(c *gin.Context) {
	c.JSON(200, gin.H{
		"category_id": c.Param("id"),
		"msg":         "implement me",
	})
	return
}

//-------------------------------
//		product routes
//-------------------------------

func (a AdminHandler) IndexProduct(c *gin.Context) {
	products, err := a.productSrv.Index(context.TODO())
	if err.Code == 404 {
		//not found
		html.Render(c, http.StatusFound, "modules/admin/html/admin_index_product", gin.H{
			"TITLE":   "index products",
			"MESSAGE": custom_error.RecordNotFound,
		})
		return
	}
	if err.Code == 500 {
		html.Error500(c)
		return
	}

	fmt.Println("found products : ", products)
	html.Render(c, http.StatusFound, "modules/admin/html/admin_index_product", gin.H{
		"TITLE":    "index products",
		"PRODUCTS": products,
	})
	return
}

func (a AdminHandler) CreateProduct(c *gin.Context) {
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
		"TITLE":      "create new product",
		"CATEGORIES": categories,
		"BRANDS":     brands.Data,
	})
	return
}

func (a AdminHandler) StoreProduct(c *gin.Context) {
	var req requests.CreateProductRequest
	_ = c.Request.ParseForm()
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("=== bind err : ", err)
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
		fmt.Println("unique error :", CheckErr)
		errors.Init()
		errors.Add("sku", custom_error.MustBeUnique)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, "/admins/products/create")
		return
	}

	fmt.Println("here======")
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

	//c.JSON(200, gin.H{
	//	"data": req,
	//})
	//return
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
	_, pErr := a.productSrv.Create(c, req)
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

func (a AdminHandler) ShowProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	selectedP, pErr := a.productSrv.Show(context.TODO(), "id", productID)
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
			"TITLE":   "show product",
			"PRODUCT": selectedP,
		},
	)
}

func (a AdminHandler) EditProduct(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	product, pErr := a.productSrv.Show(context.TODO(), "id", pID)

	if pErr.Code == 404 {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if pErr.Code == 500 {
		html.Error500(c)
		return
	}
	categories, cErr := a.categorySrv.GetAllCategories(context.TODO())
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
	fmt.Println("------- product : ", product.BrandID)
	html.Render(c, http.StatusFound, "modules/admin/html/admin_edit_product",
		gin.H{
			"TITLE":      "edit product",
			"PRODUCT":    product,
			"CATEGORIES": categories,
			"BRANDS":     brands,
		},
	)
	return
}

func (a AdminHandler) UpdateProduct(c *gin.Context) {
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
		fmt.Println("=== bind err : ", err)
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
	selectedProduct, pErr := a.productSrv.Show(context.TODO(), "id", productID)
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
	updateErr := a.productSrv.Update(c, productID, req)
	if updateErr.Code > 0 {
		fmt.Println("--- handler update product : 857 : err : ", updateErr)
		sessions.Set(c, "message", custom_messages.ProductUpdateFailed)
		c.Redirect(http.StatusFound, "/admins/products/")
	}

	sessions.Set(c, "message", custom_messages.ProductUpdatedSuccessfully)
	c.Redirect(http.StatusFound, url)
	return

}

//----------------------
//	ATTRIBUTE HANDLERS
//----------------------

func (a AdminHandler) CreateAttribute(c *gin.Context) {
	html.Render(c, http.StatusFound, "admin_create_attribute", gin.H{
		"TITLE": "create new attribute",
	})
	return
}

func (a AdminHandler) CreateAttributeValues(c *gin.Context) {

	attributes, _ := a.attributeSrv.Index(c)
	html.Render(c, http.StatusFound, "admin_create_attribute_values", gin.H{
		"TITLE":      "create new attribute-values",
		"ATTRIBUTES": attributes,
	})
	return
}

func (a AdminHandler) StoreAttribute(c *gin.Context) {
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

	newAttr, Aerr := a.attributeSrv.Create(c, req)

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

func (a AdminHandler) GetAttributesByCategoryID(c *gin.Context) {
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

func (a AdminHandler) StoreAttributeValues(c *gin.Context) {

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

	_, attErr := a.attrValueSrv.Create(c, req)

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

// ProductsAddAttributes : show html form add attributes for given product
func (a AdminHandler) ProductsAddAttributes(c *gin.Context) {

	//convert string id to int
	productID, err := strconv.Atoi(c.Param("id"))
	fmt.Println("product id  : ", productID)
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	productData, _ := a.productSrv.FetchByProductID(c, productID)

	//fetch attributes with values
	attributes, aErr := a.productSrv.FetchRootAttributes(c, productID)
	if aErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if aErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	html.Render(c, 200, "admin_create_product_att_values", gin.H{
		"TITLE":      "Add Attribute-Value to a Product",
		"ATTRIBUTES": attributes,
		"PRODUCT_ID": productID,
		"PRODUCT":    productData,
	})
	return

}

func (a AdminHandler) StoreProductsAddAttributes(c *gin.Context) {

	//convert string id to int
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	url := fmt.Sprintf("/admins/products/%d/add-attributes", productID)

	attributes := c.PostFormArray("attributes")

	//store attributes
	pAttrErr := a.productSrv.AddAttributeValues(c, productID, attributes)
	if pAttrErr.Code == 404 {
		fmt.Println("error ---------- store product attribute :", err)
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, url)
		return
	}
	if pAttrErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, url)
		return
	}

	sessions.Set(c, "message", " ویژگی ها برای محصول با موفقیت ایجاد گردید | لطفا موجودی اضافه نمایید")

	c.Redirect(http.StatusFound, fmt.Sprintf("/admins/products/%d/add-inventory", productID))
	return
}

func (a AdminHandler) ShowProductInventory(c *gin.Context) {
	productID, pErr := strconv.Atoi(c.Param("id"))
	if pErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	//fetch product data and product-attributes relation
	fmt.Println("show product inventory : product id : ", productID)
	product, err := a.productSrv.FetchByProductID(c, productID)

	if err.Code == 404 {
		url := fmt.Sprintf("/admins/products/%d/add-attributes", productID)
		sessions.Set(c, "message", custom_error.RecordNotFound)

		c.Redirect(http.StatusFound, url)
		return
	}

	productData, productErr := a.productSrv.FetchProductAttributes(c, productID)
	if productErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if productErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	html.Render(c, 200, "inventory", gin.H{
		"TITTLE":     "Product Inventory",
		"PRODUCT_ID": productID,
		"PRODUCT": gin.H{
			"Inventory": productData["inventories"],
			"Product":   productData["product"],
		},
		"PRODUCT_ATTRIBUTES": product.ProductAttributes,
	})
	return
}

func (a AdminHandler) StoreProductInventory(c *gin.Context) {
	productID, pErr := strconv.Atoi(c.Param("id"))
	urlAdminProudcts := "/admins/products"
	if pErr != nil {
		errors.Init()
		errors.SetFromErrors(pErr)
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)

		c.Redirect(http.StatusFound, urlAdminProudcts)
		return
	}

	//binding
	var req requests.CreateProductInventoryRequest
	url := fmt.Sprintf("/admins/products/%d/add-inventory", productID)
	_ = c.Request.ParseForm()
	err := c.ShouldBind(&req)
	if err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		c.Redirect(http.StatusFound, url)
		return
	}

	//insert with transaction [product inventory - product-inventory-attribute ]
	iErr := a.productSrv.CreateInventory(c, productID, req)
	if iErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, urlAdminProudcts)
		return
	}
	if iErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, urlAdminProudcts)
		return
	}

	sessions.Set(c, "message", custom_messages.ProductInventoryCreatedSuccessfully)

	c.Redirect(http.StatusFound, url)
	return
}

func (a AdminHandler) ShowProductGallery(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	product, pErr := a.productSrv.Show(context.TODO(), "id", pID)

	if pErr.Code == 404 {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if pErr.Code == 500 {
		html.Error500(c)
		return
	}
	html.Render(c, http.StatusFound, "edit-gallery-product", gin.H{
		"TITLE":   "edit gallery product",
		"PRODUCT": product,
	})
	return
}

func (a AdminHandler) DeleteProductImage(c *gin.Context) {

	imageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	image, iErr := a.productSrv.FetchImage(c, imageID)
	if iErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}
	if iErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	//remove image from disk
	imageDelErr := os.Remove(image.FullPath)

	if imageDelErr != nil {
		fmt.Println("-- remove image from disk err:", imageDelErr)
	}

	//remove image form db
	rImageErr := a.productSrv.RemoveImage(c, imageID)
	if rImageErr.Code > 0 {
		fmt.Println("---- RemoveImage --- err : ", rImageErr)
	}

	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}

func (a AdminHandler) UploadProductImages(c *gin.Context) {

	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products/")
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

		c.Redirect(http.StatusFound, c.Request.Referer())
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

			c.Redirect(http.StatusFound, c.Request.Referer())
			return

		}

		//generate file name
		imageGenerateFileName := util.GenerateFilename(image.Filename)
		imagesStoredPath = append(imagesStoredPath, imageGenerateFileName)

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

	uploadImageErr := a.productSrv.UploadImage(c, productID, imagesStoredPath)
	if uploadImageErr.Code > 0 {
		//remove images from disk
		for _, img := range imagesStoredPath {
			_ = os.Remove(viper.GetString("Upload.products") + img)
		}

		sessions.Set(c, "message", custom_error.SomethingWrongHappened)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	sessions.Set(c, "message", custom_messages.ImagesUploadedSuccessfully)
	c.Redirect(http.StatusFound, c.Request.Referer())
	return

}

func (a AdminHandler) IndexAttribute(c *gin.Context) {

	attributes, err := a.attributeSrv.Index(c)

	if err.Code == 400 {
		c.JSON(200, gin.H{
			"data": "empty",
		})
	} else if err.Code == 500 {
		html.Error500(c)
		return
	}
	html.Render(c, 200, "admin_index_attribute", gin.H{
		"TITLE":      "Index Attributes",
		"ATTRIBUTES": attributes,
	})
	return

}

func (a AdminHandler) ShowAttribute(c *gin.Context) {
	attributeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/attributes")
		return
	}

	attribute, attErr := a.attributeSrv.Show(c, attributeID)
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
			"TITLE":     "show attribute",
			"ATTRIBUTE": attribute,
		},
	)
}

func (a AdminHandler) UpdateAttribute(c *gin.Context) {

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
		fmt.Println(" ---- update attr bErr:", bErr)
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

	updateErr := a.attributeSrv.Update(c, attributeID, req)
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

func (a AdminHandler) IndexAttributeValues(c *gin.Context) {

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

func (a AdminHandler) ShowAttributeValues(c *gin.Context) {
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

func (a AdminHandler) EditAttributeValues(c *gin.Context) {
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

func (a AdminHandler) UpdateAttributeValues(c *gin.Context) {

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
	updateErr := a.attrValueSrv.Update(c, attValID, req)
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

func (a AdminHandler) DeleteProductInventoryAttribute(c *gin.Context) {
	productInventoryAttributeID, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	deleteErr := a.productSrv.DeleteInventoryAttribute(c, productInventoryAttributeID)
	if deleteErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if deleteErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	sessions.Set(c, "message", custom_messages.DeleteSuccessfully)

	//todo: when delete is successfully we have to redirect to admins/products/:id/add-inventory
	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}

func (a AdminHandler) DeleteInventory(c *gin.Context) {
	inventoryID, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	//delete
	dErr := a.productSrv.DeleteInventory(c, inventoryID)
	if dErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if dErr.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	sessions.Set(c, "message", custom_messages.DeleteSuccessfully)
	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}

func (a AdminHandler) AppendAttribute(c *gin.Context) {

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

func (a AdminHandler) UpdateQuantity(c *gin.Context) {
	//convert string id to int
	inventoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	var req = struct {
		Quantity uint `form:"update_quantity"`
	}{}

	bindErr := c.ShouldBind(&req)
	if bindErr != nil {
		sessions.Set(c, "UPDATEMESSAGE", custom_error.SomethingWrongHappened)
		c.JSON(200, gin.H{"err": bindErr})
		return
	}

	if req.Quantity == 0 {
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	updateErr := a.productSrv.UpdateInventoryQuantity(c, inventoryID, req.Quantity)
	if updateErr.Code > 0 {
		//pass appropriate msg
		sessions.Set(c, "message", custom_error.SomethingWrongHappened)
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}
	sessions.Set(c, "message", custom_error.SuccessfullyUpdated)
	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}

func (a AdminHandler) IndexCustomer(c *gin.Context) {
	customers, err := a.customerSrv.Index(c)
	if err.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
	}
	if err.Code == 500 {
		sessions.Set(c, "message", custom_error.InternalServerError)
	}

	fmt.Println("----- customers data : ", customers)
	html.Render(c, http.StatusFound, "admin_index_customer", gin.H{
		"TITLE":     "مدیریت مشتریان",
		"CUSTOMERS": customers,
	})

}
func (a AdminHandler) CreateProductFeature(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	product, pErr := a.productSrv.Show(context.TODO(), "id", pID)

	if pErr.Code > 0 {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	html.Render(c, http.StatusFound, "admin_add_product_feature",
		gin.H{
			"TITLE":   "add product feature",
			"PRODUCT": product,
		},
	)
	return
}

func (a AdminHandler) StoreProductFeature(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	url := fmt.Sprintf("/admins/products/%d/add-feature", pID)
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	_, pErr := a.productSrv.Show(c, "id", pID)
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

	addErr := a.productSrv.AddFeature(c, pID, req)
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

func (a AdminHandler) ShowProductFeature(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	//url := fmt.Sprintf("/admins/products/%d/show-feature", pID)
	if err != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	productData, pErr := a.productSrv.Show(c, "id", pID)
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

func (a AdminHandler) DeleteProductFeature(c *gin.Context) {
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

func (a AdminHandler) EditProductFeature(c *gin.Context) {
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

func (a AdminHandler) UpdateProductFeature(c *gin.Context) {
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

	updateErr := a.productSrv.UpdateFeature(c, pID, fID, req)
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
