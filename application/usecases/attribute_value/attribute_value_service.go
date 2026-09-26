package attributeValue

import (
	"context"
	"github.com/gin-gonic/gin"
	"shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
)

type AttributeValueService struct {
	repo repositories.AttributeValueRepositoryInterface
}

func NewAttributeValueService(repo repositories.AttributeValueRepositoryInterface) AttributeValueServiceInterface {
	return &AttributeValueService{repo: repo}
}

//-------------------------
//>>>>>>>> Methods <<<<<<<<
//-------------------------

func (av *AttributeValueService) Create(ctx context.Context, req *requests.CreateAttributeValueRequest) (*responses.AttributeValue, domain_err.CustomError) {

	newAttrValue, err := av.repo.Store(ctx, req)

	if err != nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}

	return responses.ToAttributeValue(newAttrValue), domain_err.CustomError{}
}

// IndexAttribute get attributes by its attribute-values relation
func (av *AttributeValueService) IndexAttribute(c *gin.Context) (*responses.Attributes, domain_err.CustomError) {
	attributes, err := av.repo.GetAllAttribute(c)
	if err != nil || attributes == nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return responses.ToAttributes(attributes), domain_err.CustomError{}
}
func (av *AttributeValueService) Show(c *gin.Context, attributeValueID int) (*responses.AttributeValue, domain_err.CustomError) {
	attValue, err := av.repo.Find(c, attributeValueID)
	if err != nil || attValue == nil {
		return nil, domain_err.HandleError(err, domain_err.RecordNotFound)
	}

	return responses.ToAttributeValue(attValue), domain_err.CustomError{}
}
func (av *AttributeValueService) Update(c *gin.Context, attributeValueID int, req *requests.UpdateAttributeValueRequest) domain_err.CustomError {
	_, err := av.repo.Update(c, attributeValueID, req)

	if err != nil {
		return domain_err.HandleError(err, domain_err.RecordNotFound)
	}
	return domain_err.CustomError{}

}
