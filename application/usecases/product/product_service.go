package product

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
	"log"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
	"strings"
)

type ProductService struct {
	repo repositories.ProductRepositoryInterface
}

func NewProductService(repo repositories.ProductRepositoryInterface) ProductServiceInterface {
	return &ProductService{repo: repo}
}

//-----------------------------------------
//<<<<<<<<<<<<<<<< Method >>>>>>>>>>>>>>>>>
//-----------------------------------------

func (p *ProductService) Index(ctx context.Context) (*responses.Products, domain_err.CustomError) {

	products, err := p.repo.GetAll(ctx)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToProducts(products), domain_err.CustomError{}
}

func (p *ProductService) Show(ctx context.Context, columnName string, value any) (*responses.Product, []bson.M, domain_err.CustomError) {

	pResult, err := p.repo.FindBy(ctx, columnName, value)
	if err != nil {
		return nil, nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}

	mongoProduct, err := p.repo.GetAllMongoProduct(ctx)
	log.Println("--- get all mongo product :", len(mongoProduct), " | err:", err)
	if err != nil {
		return nil, nil, domain_err.CustomError{}
	} else {
		//add id field which is string not object
		for i := range mongoProduct {
			if objID, ok := mongoProduct[i]["_id"].(primitive.ObjectID); ok {
				mongoProduct[i]["id"] = objID.Hex()
			}
		}
	}

	return responses.ToProduct(pResult), mongoProduct, domain_err.CustomError{}
}

func (p *ProductService) Create(ctx context.Context, req *requests.CreateProductRequest) (*responses.Product, domain_err.CustomError) {

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

	newProduct, err := p.repo.Store(ctx, &prepareProduct)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToProduct(newProduct), domain_err.CustomError{}
}

func prepareProductImages(imageNames []string) []*entities.ProductImages {
	var pImages []*entities.ProductImages
	for _, imageName := range imageNames {
		pImages = append(pImages, &entities.ProductImages{Path: imageName})
	}
	return pImages
}

func (p *ProductService) CheckSkuIsUnique(ctx context.Context, sku string) (bool, domain_err.CustomError) {
	_, err := p.repo.FindBy(ctx, "sku", sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, domain_err.New(err.Error(), domain_err.RecordNotFound, 404)
		}
		return false, domain_err.New(err.Error(), domain_err.InternalServerError, 500)
	}
	return false, domain_err.CustomError{}
}

func (p *ProductService) FetchByProductID(c *gin.Context, productID int) (*responses.Product, domain_err.CustomError) {
	pResult, err := p.repo.FindByID(c, productID)

	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToProduct(pResult), domain_err.CustomError{}
}

func (p *ProductService) FetchRootAttributes(c *gin.Context, productID int) (*responses.Attributes, domain_err.CustomError) {

	attributes, err := p.repo.GetRootAttributes(c, productID)

	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToAttributes(attributes), domain_err.CustomError{}
}

func (p *ProductService) AddAttributeValues(c *gin.Context, productID int, attributes []string) domain_err.CustomError {

	if err := p.repo.StoreAttributeValues(c, productID, attributes); err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}

	return domain_err.CustomError{}
}

func (p *ProductService) FetchProductAttributes(c *gin.Context, productID int) (map[string]interface{}, domain_err.CustomError) {
	//fetch product and its attribute and also inventories
	pResult, err := p.repo.GetProductAndAttributes(c, productID)
	if err != nil {
		return pResult, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return pResult, domain_err.CustomError{}
}

func (p *ProductService) CreateInventory(c *gin.Context, productID int, req *requests.CreateProductInventoryRequest) domain_err.CustomError {
	_, err := p.repo.StoreProductInventory(c, productID, req)
	if err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

func (p *ProductService) FetchImage(c *gin.Context, imageID int) (*responses.ImageProduct, domain_err.CustomError) {
	image, err := p.repo.GetImage(c, imageID)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToImageProduct(image), domain_err.CustomError{}
}

func (p *ProductService) RemoveImage(c *gin.Context, imageID int) domain_err.CustomError {
	if err := p.repo.DeleteImage(c, imageID); err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

func (p *ProductService) UploadImage(c *gin.Context, productID int, imageStoredPath []string) domain_err.CustomError {
	if err := p.repo.StoreImages(c, productID, imageStoredPath); err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

func (p *ProductService) Update(c *gin.Context, productID int, req *requests.UpdateProductRequest) domain_err.CustomError {
	_, err := p.repo.Update(c, productID, req)
	if err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

// DeleteInventoryAttribute delete record form product_inventory_attributes table
func (p *ProductService) DeleteInventoryAttribute(c *gin.Context, productInventoryAttributeID int) domain_err.CustomError {
	err := p.repo.DeleteInventoryAttribute(c, productInventoryAttributeID)
	if err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

func (p *ProductService) DeleteInventory(c *gin.Context, inventoryID int) domain_err.CustomError {
	if err := p.repo.DeleteInventory(c, inventoryID); err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

func (p *ProductService) AppendAttributesToInventory(c *gin.Context, inventoryID int, attributes []string) domain_err.CustomError {
	if err := p.repo.AppendAttributesToInventory(c, inventoryID, attributes); err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}
func (p *ProductService) UpdateInventoryQuantity(c *gin.Context, inventoryID int, quantity uint) domain_err.CustomError {
	if err := p.repo.UpdateInventoryQuantity(c, inventoryID, quantity); err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

func (p *ProductService) AddFeature(c *gin.Context, productID int, req *requests.CreateProductFeatureRequest) domain_err.CustomError {
	err := p.repo.InsertFeature(c, productID, req)
	if err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

func (p *ProductService) RemoveFeature(c *gin.Context, productID int, featureID int) domain_err.CustomError {
	err := p.repo.DeleteFeature(c, productID, featureID)
	if err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

func (p *ProductService) FetchFeature(c *gin.Context, productID int, featureID int) (*responses.Feature, domain_err.CustomError) {
	feat, err := p.repo.GetFeatureBy(c, productID, featureID)

	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}

	return responses.ToFeature(feat), domain_err.CustomError{}
}

func (p *ProductService) UpdateFeature(c *gin.Context, productID int, featureID int, req *requests.UpdateProductFeatureRequest) domain_err.CustomError {
	if err := p.repo.EditFeature(c, productID, featureID, req); err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}
func (p *ProductService) AddRecommendation(c *gin.Context, productID int, productRecommendationIDs []string) domain_err.CustomError {
	err := p.repo.InsertRecommendation(c, productID, productRecommendationIDs)
	if err != nil {
		log.Println("--- AddRecommendation err: ", err)
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}
}

func (p *ProductService) FetchAllRecommendation(c *gin.Context, productID int) ([]bson.M, domain_err.CustomError) {
	recommendations, err := p.repo.GetAllRecommendation(c, productID)
	log.Println("--- fetch all recommendations : ", len(recommendations), " | err:", err)
	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return recommendations, domain_err.CustomError{}
}
