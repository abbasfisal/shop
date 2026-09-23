package repositories

import (
	"context"
	"shop/domain/entities"
)

type AuthenticateRepositoryInterface interface {
	FindBy(ctx context.Context, phone string) (*entities.User, error)
	FindByUserID(ctx context.Context, userID uint) (*entities.User, error)
}
