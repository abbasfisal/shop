package handlers

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
	"shop/internal/modules/admin/services/auth"
	"shop/internal/modules/admin/services/category"
	"shop/internal/modules/admin/services/product"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/old"
	"shop/internal/pkg/sessions"
	"shop/internal/pkg/util"
	"slices"
	"strconv"
)

type AdminHandler struct {
	authSrv     auth.AuthenticateServiceInterface
	categorySrv category.CategoryServiceInterface
	productSrv  product.ProductServiceInterface
	i18nBundle  *i18n.Bundle
	//order service
	//user service
	//cart service
}

func NewAdminHandler(
	authSrv auth.AuthenticateServiceInterface,
	categorySrv category.CategoryService,
	productSrv product.ProductServiceInterface,
	i18nBundle *i18n.Bundle,
) AdminHandler {
	return AdminHandler{
		authSrv:     authSrv,
		categorySrv: categorySrv,
		productSrv:  productSrv,
		i18nBundle:  i18nBundle,
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
	catID, err := strconv.Atoi(c.Param("00"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/categories/")
		return
	}

	cat, cErr := a.categorySrv.Show(context.TODO(), catID)

	if cErr.Code == 500 {
		html.Error500(c)
		return
	}
	if cErr.Code > 0 {
		c.Redirect(http.StatusFound, "/admins/categories")
		return
	}

	html.Render(c, http.StatusFound, "modules/admin/html/admin_edit_category", gin.H{
		"TITLE":    "edit category",
		"CATEGORY": cat,
	})
	return
}

func (a AdminHandler) UpdateCategory(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "implement me",
	})
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
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if err.Code == 500 {
		html.Error500(c)
		return
	}
	html.Render(c, http.StatusFound, "modules/admin/html/admin_create_product", gin.H{
		"TITLE":      "create new product",
		"CATEGORIES": categories,
	})
	return
}

func (a AdminHandler) StoreProduct(c *gin.Context) {
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

	//check category_id existence
	_, cErr := a.categorySrv.Show(c, req.CategoryID)
	if cErr.Code > 0 {
		html.Error500(c)
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
		html.Error500(c)
		return
	}

	html.Render(c, http.StatusFound, "modules/admin/html/admin_show_product",
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

	html.Render(c, http.StatusFound, "modules/admin/html/admin_edit_product",
		gin.H{
			"TITLE":      "edit product",
			"PRODUCT":    product,
			"CATEGORIES": categories,
		},
	)
	return
}

func (a AdminHandler) UpdateProduct(c *gin.Context) {}

//----------------------
//	ATTRIBUTE HANDLERS
//----------------------

func (a AdminHandler) CreateAttribute(c *gin.Context) {
	categories, _ := a.categorySrv.GetAllParentCategory(c)
	html.Render(c, http.StatusFound, "admin_create_attribute", gin.H{
		"TITLE":      "create new attribute",
		"CATEGORIES": categories,
	})
	return
}

func (a AdminHandler) CreateAttributeValues(c *gin.Context) {

	//category where parent_id=null
	//attribute will fetch by ajax request
	categories, _ := a.categorySrv.GetAllParentCategory(c)
	html.Render(c, http.StatusFound, "admin_create_attribute_values", gin.H{
		"TITLE":      "create new attribute-values",
		"CATEGORIES": categories,
		"ATTRIBUTES": "attributes",
	})
	return
}
