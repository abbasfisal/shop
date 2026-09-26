package attribute

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/domain/entities"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
)

type AttributeService struct {
	repo repositories.AttributeRepositoryInterface
}

func NewAttributeService(repo repositories.AttributeRepositoryInterface) AttributeServiceInterface {
	return &AttributeService{repo: repo}
}

//-----------------
//>>>> Methods <<<<
//-----------------

func (as AttributeService) Create(ctx context.Context, req *requests.CreateAttributeRequest) (*responses.Attribute, error) {

	attr := entities.Attribute{
		Title:     req.Title,
		InputType: req.NormalizedInputType(),
	}

	result, err := as.repo.Store(ctx, &attr)
	if err != nil || result == nil {
		return nil, err
	}
	return responses.ToAttribute(result), nil
}
func (as AttributeService) FetchByCategoryID(ctx context.Context, categoryID int) (*responses.Attributes, error) {
	attributes, err := as.repo.GetByCategory(ctx, categoryID)
	return responses.ToAttributes(attributes), err
}
func (as AttributeService) Index(c *gin.Context) (*responses.Attributes, domain_err.CustomError) {
	attributes, err := as.repo.GetAll(c)
	if err != nil || attributes == nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToAttributes(attributes), domain_err.CustomError{}
}
func (as AttributeService) Show(c context.Context, attributeID int) (*responses.Attribute, domain_err.CustomError) {
	att, err := as.repo.GetByID(c, attributeID)
	if err != nil || att == nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToAttribute(att), domain_err.CustomError{}
}
func (as AttributeService) Update(c *gin.Context, attributeID int, req *requests.CreateAttributeRequest) domain_err.CustomError {
	err := as.repo.Update(c, attributeID, req)
	if err != nil {
		return domain_err.HandleError(err, domain_err.IDIsNotCorrect)
	}
	return domain_err.CustomError{}
}
