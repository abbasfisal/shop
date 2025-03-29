package handlers

import "C"
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/services/attribute"
	attributeValue "shop/internal/modules/admin/services/attribute_value"
	"shop/internal/modules/admin/services/auth"
	"shop/internal/modules/admin/services/banner"
	"shop/internal/modules/admin/services/brand"
	"shop/internal/modules/admin/services/category"
	"shop/internal/modules/admin/services/customer"
	"shop/internal/modules/admin/services/dashboard"
	"shop/internal/modules/admin/services/order"
	"shop/internal/modules/admin/services/product"
	"shop/internal/pkg/bootstrap"
	"shop/internal/pkg/custom_error"
	"shop/internal/pkg/custom_messages"
	"shop/internal/pkg/errors"
	"shop/internal/pkg/html"
	"shop/internal/pkg/old"
	"shop/internal/pkg/sessions"
	"shop/internal/pkg/util"
	"slices"
	"strconv"
)

type AdminHandler struct {
	authSrv          auth.AuthenticateServiceInterface
	categorySrv      category.CategoryServiceInterface
	productSrv       product.ProductServiceInterface
	attributeSrv     attribute.AttributeServiceInterface
	attrValueSrv     attributeValue.AttributeValueServiceInterface
	brandSrv         brand.BrandServiceInterface
	customerSrv      customer.CustomerServiceInterface
	orderSrv         order.OrderServiceInterface
	DashboardService *dashboard.DashboardService
	bannerSrv        *banner.BannerService

	dep *bootstrap.Dependencies
}

func NewAdminHandler(
	authSrv auth.AuthenticateServiceInterface,
	categorySrv category.CategoryServiceInterface,
	productSrv product.ProductServiceInterface,
	attributeSrv attribute.AttributeServiceInterface,
	attrValueSrv attributeValue.AttributeValueServiceInterface,
	brandSrv brand.BrandServiceInterface,
	customerSrv customer.CustomerServiceInterface,
	orderSrv order.OrderServiceInterface,
	dashboardSrv *dashboard.DashboardService,
	bannerSrv *banner.BannerService,

	dep *bootstrap.Dependencies,
) *AdminHandler {
	return &AdminHandler{
		authSrv:          authSrv,
		categorySrv:      categorySrv,
		productSrv:       productSrv,
		attributeSrv:     attributeSrv,
		attrValueSrv:     attrValueSrv,
		brandSrv:         brandSrv,
		customerSrv:      customerSrv,
		orderSrv:         orderSrv,
		DashboardService: dashboardSrv,
		bannerSrv:        bannerSrv,

		dep: dep,
	}
}

func (a *AdminHandler) ShowHome(c *gin.Context) {
	data, err := a.DashboardService.GetStaticalData()
	if err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	html.Render(c, http.StatusOK, "modules/admin/html/admin_home", gin.H{
		"TITLE": "داشبورد",
		"DATA":  data,
	})
	return
}

// ProductsAddAttributes : show html form add attributes for given product
func (a *AdminHandler) ProductsAddAttributes(c *gin.Context) {

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

func (a *AdminHandler) StoreProductsAddAttributes(c *gin.Context) {

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

func (a *AdminHandler) ShowProductInventory(c *gin.Context) {
	productID, pErr := strconv.Atoi(c.Param("id"))
	if pErr != nil {
		sessions.Set(c, "message", custom_error.IDIsNotCorrect)
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}

	//fetch productResult data and productResult-attributes relation
	fmt.Println("show productResult inventory : productResult id : ", productID)
	productResult, err := a.productSrv.FetchByProductID(c, productID)

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
			"Product":   productData["productResult"],
		},
		"PRODUCT_ATTRIBUTES": productResult.ProductAttributes,
	})
	return
}

func (a *AdminHandler) StoreProductInventory(c *gin.Context) {
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
	iErr := a.productSrv.CreateInventory(c, productID, &req)
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

func (a *AdminHandler) ShowProductGallery(c *gin.Context) {
	pID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admins/products/")
		return
	}

	productShow, _, pErr := a.productSrv.Show(c.Request.Context(), "id", pID)

	if pErr.Code == 404 {
		c.Redirect(http.StatusFound, "/admins/products")
		return
	}
	if pErr.Code == 500 {
		html.Error500(c)
		return
	}
	html.Render(c, http.StatusFound, "edit-gallery-product", gin.H{
		"TITLE":      "ویرایش تصاویر محصول",
		"PRODUCT":    productShow,
		"MEDIA_PATH": util.GetProductStoragePath(),
	})
	return
}

func (a *AdminHandler) DeleteProductImage(c *gin.Context) {

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

	//storage s3
	if os.Getenv("STORAGE_STATUS") == "active" {
		go func() {
			fmt.Println("-- image address to delete : ", image.OriginalPath)
			err := a.dep.Storage.DeleteFile(os.Getenv("STORAGE_PRODUCT_PATH") + image.OriginalPath)
			if err != nil {
				log.Println("-- s3 delete file failed: ", err)
			}
		}()
	}

	//remove image form db
	rImageErr := a.productSrv.RemoveImage(c, imageID)
	if rImageErr.Code > 0 {
		fmt.Println("---- RemoveImage --- err : ", rImageErr)
	}

	c.Redirect(http.StatusFound, c.Request.Referer())
	return
}

func (a *AdminHandler) UploadProductImages(c *gin.Context) {

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

		//upload in s3
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

func (a *AdminHandler) DeleteProductInventoryAttribute(c *gin.Context) {
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

func (a *AdminHandler) DeleteInventory(c *gin.Context) {
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

func (a *AdminHandler) UpdateQuantity(c *gin.Context) {
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
