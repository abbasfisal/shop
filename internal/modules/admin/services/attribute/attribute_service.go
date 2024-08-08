package attribute

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/repositories/attribute"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type AttributeService struct {
	repo attribute.AttributeRepositoryInterface
}

func NewAttributeService(repo attribute.AttributeRepositoryInterface) AttributeService {
	return AttributeService{repo: repo}
}

func (as AttributeService) Create(ctx context.Context, req requests.CreateAttributeRequest) (responses.Attribute, error) {
	var res responses.Attribute

	attr := entities.Attribute{
		//CategoryID: req.CategoryID,
		Title: req.Title,
	}

	result, err := as.repo.Store(ctx, attr)
	if err != nil {
		return res, err
	}
	return responses.ToAttribute(result), nil
}

func (as AttributeService) FetchByCategoryID(ctx context.Context, categoryID int) (responses.Attributes, error) {
	attributes, err := as.repo.GetByCategory(ctx, categoryID)
	return responses.ToAttributes(attributes), err
}
func (as AttributeService) Index(c *gin.Context) (responses.Attributes, custom_error.CustomError) {
	attributes, err := as.repo.GetAll(c)
	if err != nil {
		return responses.Attributes{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToAttributes(attributes), custom_error.CustomError{}
}
