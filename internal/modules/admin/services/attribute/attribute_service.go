package attribute

import (
	"context"
	"shop/internal/entities"
	"shop/internal/modules/admin/repositories/attribute"
	"shop/internal/modules/admin/requests"
	"shop/internal/modules/admin/responses"
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
		CategoryID: req.CategoryID,
		Title:      req.Title,
	}

	result, err := as.repo.Store(ctx, attr)
	if err != nil {
		return res, err
	}
	return responses.ToAttribute(result), nil
}
