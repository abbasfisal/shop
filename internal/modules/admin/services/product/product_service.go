package product

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"shop/internal/entities"
	"shop/internal/modules/admin/repositories/product"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
	"strings"
)

type ProductService struct {
	repo product.ProductRepositoryInterface
}

func NewProductService(repo product.ProductRepositoryInterface) *ProductService {
	return &ProductService{repo: repo}
}

func (p ProductService) Index(ctx context.Context) (responses.Products, custom_error.CustomError) {
	var response responses.Products
	products, err := p.repo.GetAll(ctx)
	if err != nil {
		return response, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProducts(products), custom_error.CustomError{}
}

func (p ProductService) Show(ctx context.Context, columnName string, value any) (responses.Product, custom_error.CustomError) {

	product, err := p.repo.FindBy(ctx, columnName, value)
	if err != nil {
		return responses.Product{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProduct(product), custom_error.CustomError{}
}

func (p ProductService) Create(ctx context.Context, req requests.CreateProductRequest) (responses.Product, custom_error.CustomError) {

	var prepareProduct = entities.Product{
		CategoryID: uint(req.CategoryID),
		BrandID:    req.BrandID,
		Title:      strings.TrimSpace(req.Title),
		Slug:       strings.TrimSpace(req.Title),
		Sku:        strings.TrimSpace(req.Title),
		Status: func() bool {
			if req.Status == "" {
				return false
			}
			return true
		}(),
		OriginalPrice: req.OriginalPrice,
		SalePrice:     req.SalePrice,
		Description:   strings.TrimSpace(req.Description),
		ProductImages: prepareProductImages(req.ProductImage),
	}

	newProduct, err := p.repo.Store(ctx, prepareProduct)
	if err != nil {
		return responses.Product{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProduct(newProduct), custom_error.CustomError{}
}

func prepareProductImages(imageNames []string) []entities.ProductImages {
	var pImages []entities.ProductImages
	for _, imageName := range imageNames {
		pImages = append(pImages, entities.ProductImages{Path: imageName})
	}
	return pImages
}

func (p ProductService) CheckSkuIsUnique(ctx context.Context, sku string) (bool, custom_error.CustomError) {
	_, err := p.repo.FindBy(ctx, "sku", sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return false, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}
	return false, custom_error.CustomError{}
}

func (p ProductService) FetchByProductID(c *gin.Context, productID int) (responses.Product, custom_error.CustomError) {
	product, err := p.repo.FindByID(c, productID)

	if err != nil {
		return responses.Product{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToProduct(product), custom_error.CustomError{}
}

func (p ProductService) FetchRootAttributes(c *gin.Context, productID int) (responses.Attributes, custom_error.CustomError) {

	attributes, err := p.repo.GetRootAttributes(c, productID)

	if err != nil {
		return responses.Attributes{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToAttributes(attributes), custom_error.CustomError{}
}

func (p ProductService) AddAttributeValues(c *gin.Context, productID int, attributes []string) custom_error.CustomError {

	if err := p.repo.StoreAttributeValues(c, productID, attributes); err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}

	return custom_error.CustomError{}
}

func (p ProductService) FetchProductAttributes(c *gin.Context, productID int) (map[string]interface{}, custom_error.CustomError) {
	//fetch product and its attribute and also inventories
	product, err := p.repo.GetProductAndAttributes(c, productID)
	if err != nil {
		return product, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return product, custom_error.CustomError{}
}

func (p ProductService) CreateInventory(c *gin.Context, productID int, req requests.CreateProductInventoryRequest) custom_error.CustomError {
	_, err := p.repo.StoreProductInventory(c, productID, req)
	if err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return custom_error.CustomError{}
}

func (p ProductService) FetchImage(c *gin.Context, imageID int) (responses.ImageProduct, custom_error.CustomError) {
	image, err := p.repo.GetImage(c, imageID)
	if err != nil {
		return responses.ImageProduct{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToImageProduct(image), custom_error.CustomError{}
}

func (p ProductService) RemoveImage(c *gin.Context, imageID int) custom_error.CustomError {
	if err := p.repo.DeleteImage(c, imageID); err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return custom_error.CustomError{}
}

func (p ProductService) UploadImage(c *gin.Context, productID int, imageStoredPath []string) custom_error.CustomError {
	if err := p.repo.StoreImages(c, productID, imageStoredPath); err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return custom_error.CustomError{}
}

func (p ProductService) Update(c *gin.Context, productID int, req requests.UpdateProductRequest) custom_error.CustomError {
	_, err := p.repo.Update(c, productID, req)
	if err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return custom_error.CustomError{}
}

// DeleteInventoryAttribute delete record form product_inventory_attributes table
func (p ProductService) DeleteInventoryAttribute(c *gin.Context, productInventoryAttributeID int) custom_error.CustomError {
	err := p.repo.DeleteInventoryAttribute(c, productInventoryAttributeID)
	if err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return custom_error.CustomError{}
}

func (p ProductService) DeleteInventory(c *gin.Context, inventoryID int) custom_error.CustomError {
	if err := p.repo.DeleteInventory(c, inventoryID); err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return custom_error.CustomError{}
}

func (p ProductService) AppendAttributesToInventory(c *gin.Context, inventoryID int, attributes []string) custom_error.CustomError {
	if err := p.repo.AppendAttributesToInventory(c, inventoryID, attributes); err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return custom_error.CustomError{}
}
func (p ProductService) UpdateInventoryQuantity(c *gin.Context, inventoryID int, quantity uint) custom_error.CustomError {
	if err := p.repo.UpdateInventoryQuantity(c, inventoryID, quantity); err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return custom_error.CustomError{}

}
