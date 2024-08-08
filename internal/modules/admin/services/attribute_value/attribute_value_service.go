package attributeValue

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"shop/internal/entities"
	"shop/internal/modules/admin/repositories/attribute_value"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type AttributeValueService struct {
	repo attributeValue.AttributeValueRepositoryInterface
}

func NewAttributeValueService(repo attributeValue.AttributeValueRepositoryInterface) AttributeValueService {
	return AttributeValueService{repo: repo}
}

func (av AttributeValueService) Create(ctx context.Context, req requests.CreateAttributeValueRequest) (responses.AttributeValue, custom_error.CustomError) {

	newAttrValue, err := av.repo.Store(ctx, entities.AttributeValue{
		AttributeID: req.AttributeID,
		Value:       req.Value,
	})

	if err != nil {
		return responses.AttributeValue{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}

	return responses.ToAttributeValue(newAttrValue), custom_error.CustomError{}
}

// IndexAttribute get attributes by its attribute-values relation
func (av AttributeValueService) IndexAttribute(c *gin.Context) (responses.Attributes, custom_error.CustomError) {
	attributes, err := av.repo.GetAllAttribute(c)
	if err != nil {
		return responses.Attributes{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return responses.ToAttributes(attributes), custom_error.CustomError{}
}

func (av AttributeValueService) Show(c *gin.Context, attributeValueID int) (responses.AttributeValue, custom_error.CustomError) {
	attValue, err := av.repo.Find(c, attributeValueID)
	if err != nil {
		return responses.AttributeValue{}, custom_error.HandleError(err, custom_error.RecordNotFound)
	}

	return responses.ToAttributeValue(attValue), custom_error.CustomError{}
}

func (av AttributeValueService) Update(c *gin.Context, attributeValueID int, req requests.UpdateAttributeValueRequest) custom_error.CustomError {
	attributeValue, err := av.repo.Update(c, attributeValueID, req)

	fmt.Println("--- update succ att-value : ", attributeValue)
	if err != nil {
		return custom_error.HandleError(err, custom_error.RecordNotFound)
	}
	return custom_error.CustomError{}

}
