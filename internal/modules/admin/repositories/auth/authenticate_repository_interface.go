package auth

import (
	"context"
	"shop/internal/entities"
)

type AuthenticateRepositoryInterface interface {
	FindBy(ctx context.Context, phone string) (*entities.User, error)
	FindByUserID(ctx context.Context, userID uint) (*entities.User, error)
}
