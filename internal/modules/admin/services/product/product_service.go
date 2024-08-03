package product

import (
	"context"
	"errors"
	"fmt"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return response, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}

	return responses.ToProducts(products), custom_error.CustomError{}
}

func (p ProductService) Show(ctx context.Context, columnName string, value any) (responses.Product, custom_error.CustomError) {

	product, err := p.repo.FindBy(ctx, columnName, value)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.Product{}, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return responses.Product{}, custom_error.New(err.Error(), custom_error.RecordNotFound, 500)
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
		//Quantity:      req.Quantity,
		OriginalPrice: req.OriginalPrice,
		SalePrice:     req.SalePrice,
		Description:   strings.TrimSpace(req.Description),
		ProductImages: func() []entities.ProductImages {
			var pImages []entities.ProductImages
			for _, imageName := range req.ProductImage {
				pImages = append(pImages, entities.ProductImages{
					Model:     gorm.Model{},
					ProductID: 0,
					Path:      imageName,
				})
			}
			return pImages
		}(),
	}

	newProduct, err := p.repo.Store(ctx, prepareProduct)
	if err != nil {
		return responses.Product{}, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}
	return responses.ToProduct(newProduct), custom_error.CustomError{}
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.Product{}, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return responses.Product{}, custom_error.New(err.Error(), custom_error.RecordNotFound, 500)
	}
	return responses.ToProduct(product), custom_error.CustomError{}
}

func (p ProductService) FetchRootAttributes(c *gin.Context, productID int) (responses.Attributes, custom_error.CustomError) {

	attributes, err := p.repo.GetRootAttributes(c, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.Attributes{}, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return responses.Attributes{}, custom_error.New(err.Error(), custom_error.RecordNotFound, 500)
	}

	return responses.ToAttributes(attributes), custom_error.CustomError{}

}

func (p ProductService) AddAttributeValues(c *gin.Context, productID int, attributes []string) custom_error.CustomError {
	err := p.repo.StoreAttributeValues(c, productID, attributes)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return custom_error.New(err.Error(), custom_error.InternalServerError, 500)

	}

	return custom_error.CustomError{}
}

func (p ProductService) FetchProductAttributes(c *gin.Context, productID int) (responses.Product, custom_error.CustomError) {
	//fetch product , and fetch product_attributes
	product, err := p.repo.GetProductAndAttributes(c, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.Product{}, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return responses.Product{}, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}
	return responses.ToProduct(product), custom_error.CustomError{}
}

func (p ProductService) CreateInventory(c *gin.Context, productID int, req requests.CreateProductInventoryRequest) custom_error.CustomError {
	_, err := p.repo.StoreProductInventory(c, productID, req)
	if err != nil {
		fmt.Println("create Inventory Error : ", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}
	return custom_error.CustomError{}
}

func (p ProductService) FetchImage(c *gin.Context, imageID int) (responses.ImageProduct, custom_error.CustomError) {
	image, err := p.repo.GetImage(c, imageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.ImageProduct{}, custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return responses.ImageProduct{}, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}
	return responses.ToImageProduct(image), custom_error.CustomError{}
}

func (p ProductService) RemoveImage(c *gin.Context, imageID int) custom_error.CustomError {
	err := p.repo.DeleteImage(c, imageID)
	if err != nil {
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
			}
			return custom_error.New(err.Error(), custom_error.InternalServerError, 500)
		}
	}
	return custom_error.CustomError{}
}

func (p ProductService) UploadImage(c *gin.Context, productID int, imageStoredPath []string) custom_error.CustomError {
	err := p.repo.StoreImages(c, productID, imageStoredPath)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return custom_error.New(err.Error(), custom_error.RecordNotFound, 404)
		}
		return custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}
	return custom_error.CustomError{}
}
