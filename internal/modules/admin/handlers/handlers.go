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
	"shop/internal/modules/admin/services/attribute"
	attributeValue "shop/internal/modules/admin/services/attribute_value"
	"shop/internal/modules/admin/services/auth"
	"shop/internal/modules/admin/services/brand"
	"shop/internal/modules/admin/services/category"
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
	i18nBundle *i18n.Bundle,
) AdminHandler {
	return AdminHandler{
		authSrv:      authSrv,
		categorySrv:  categorySrv,
		productSrv:   productSrv,
		attributeSrv: attributeSrv,
		attrValueSrv: attrValueSrv,
		brandSrv:     brandSrv,

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
		fmt.Println("----- string to id err : ", catIDErr)
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
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if err.Code == 500 {
		html.Error500(c)
		return
	}
	brands, _ := a.brandSrv.Index(c)
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
	//c.JSON(200, gin.H{
	//	"products:": selectedP,
	//})
	//return
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

func (a AdminHandler) UpdateProduct(c *gin.Context) {}

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
			errors.Add("attribute_id", "The field is required.")
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

	newAttrValue, err := a.attrValueSrv.Create(c, req)

	if err != nil || newAttrValue.ID <= 0 {
		sessions.Set(c, "message", "خطا در ایجاد ویژگی")
		c.Redirect(http.StatusFound, "/admins/attribute-values/create")
		return
	}

	sessions.Set(c, "message", "ایجاد ویژگی با موفقیت انجام شد")
	c.Redirect(http.StatusFound, "/admins/attribute-values/create")
	return
}

// ProductsAddAttributes : show html form add attributes for given product
func (a AdminHandler) ProductsAddAttributes(c *gin.Context) {

	//fetch product
	productID, err := strconv.Atoi(c.Param("id"))
	fmt.Println("product id  : ", productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
		return
	}

	//fetch attributes with values
	attributes, aErr := a.productSrv.FetchRootAttributes(c, productID)

	if aErr.Code == 404 {
		sessions.Set(c, "message", custom_error.RecordNotFound)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if aErr.Code == 500 {
		html.Error500(c)
		return
	}
	// end fetch attributes

	html.Render(c, 200, "att", gin.H{
		"ATTRIBUTES": attributes,
		"PRODUCT_ID": productID,
	})
	return

}

func (a AdminHandler) StoreProductsAddAttributes(c *gin.Context) {

	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(422, gin.H{
			"error": "fail converting id to integer",
		})
		return
	}

	attributes := c.PostFormArray("attributes")

	pAttrErr := a.productSrv.AddAttributeValues(c, productID, attributes)
	if pAttrErr.Code == 404 {
		c.JSON(404, gin.H{
			"message":      "product not found",
			"err":          pAttrErr.DisplayMessage,
			"err original": pAttrErr.OriginalMessage,
		})
		return
	}
	if pAttrErr.Code == 500 {
		html.Error500(c)

		return
	}

	sessions.Set(c, "message", " ویژگی ها برای محصول با موفقیت ایجاد گردید | لطفا موجودی اضافه نمایید")
	url := fmt.Sprintf("/admins/products/%d/add-inventory", productID)
	c.Redirect(http.StatusFound, url)
	return
}

func (a AdminHandler) ShowProductInventory(c *gin.Context) {
	productID, pErr := strconv.Atoi(c.Param("id"))
	if pErr != nil {
		c.JSON(429, gin.H{
			"error": pErr.Error(),
		})
		return
	}
	//fetch product data
	fmt.Println("show product inventory : product id : ", productID)
	//product, err := a.productSrv.FetchByProductID(c, productID)
	product, err := a.productSrv.FetchProductAttributes(c, productID)

	if err.Code == 404 || len(product.ProductAttributes.Data) <= 0 {
		url := fmt.Sprintf("/admins/products/%d/add-attributes", productID)

		sessions.Set(c, "message", "ابتدا ویژگی برای محصول مورد نظر اضافه کنید")
		c.Redirect(http.StatusFound, url)
		return
	}
	//c.JSON(200, gin.H{
	//	"data": product,
	//})
	//return
	html.Render(c, 200, "inventory", gin.H{
		"TITTLE":     "Product Inventory",
		"PRODUCT_ID": product.ID,
		"PRODUCT":    product,
	})
	return
}

func (a AdminHandler) StoreProductInventory(c *gin.Context) {
	productID, pErr := strconv.Atoi(c.Param("id"))
	if pErr != nil {
		errors.Init()
		errors.SetFromErrors(pErr)
		sessions.Set(c, "message", "id محصول اشتباه است")

		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	//validation
	var req requests.CreateProductInventoryRequest
	c.Request.ParseForm()
	_ = c.Request.ParseForm()
	err := c.ShouldBind(&req)
	if err != nil {
		errors.Init()
		errors.SetFromErrors(err)
		sessions.Set(c, "errors", errors.ToString())

		old.Init()
		old.Set(c)
		sessions.Set(c, "olds", old.ToString())

		url := fmt.Sprintf("/admins/products/%d/add-inventory", productID)
		c.Redirect(http.StatusFound, url)
		return
	}
	//end validation

	iErr := a.productSrv.CreateInventory(c, productID, req)
	if iErr.Code == 404 {
		sessions.Set(c, "message", "رکورد یافت نشد")
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if iErr.Code == 500 {
		html.Error500(c)
		return
	}
	fmt.Println("inventory created")
	sessions.Set(c, "message", "موجودی برای محصول با موفقیت اضافه گردید")
	url := fmt.Sprintf("/admins/products/%d/add-inventory", productID)
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

	if imageDelErr == nil {
		//remove image form db
		rImageErr := a.productSrv.RemoveImage(c, imageID)

		if rImageErr.Code > 0 {
			//image record  removed from db
			fmt.Println("---- RemoveImage --- err : ", rImageErr)
			if rImageErr.Code == 404 {
				sessions.Set(c, "message", custom_error.RecordNotFound)
				c.Redirect(http.StatusFound, c.Request.Referer())
				return
			}
			if rImageErr.Code == 500 {
				sessions.Set(c, "message", custom_error.InternalServerError)
				c.Redirect(http.StatusFound, c.Request.Referer())
				return
			}
		} else {
			// image record deleted successfully
			sessions.Set(c, "message", custom_messages.ImageDeletedSuccessfully)
			c.Redirect(http.StatusFound, c.Request.Referer())
		}
	}

	fmt.Println("---- image not removed form disk : err: ------ ", imageDelErr)
	sessions.Set(c, "message", custom_messages.ImageNotRemovedFromDisk)
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
