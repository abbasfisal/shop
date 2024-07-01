package auth

import (
	"context"
	"gorm.io/gorm"
	"shop/internal/entities"
)

type AuthenticateRepository struct {
	db *gorm.DB
}

func NewAuthenticateRepository(db *gorm.DB) AuthenticateRepository {
	return AuthenticateRepository{
		db: db,
	}
}

func (a AuthenticateRepository) FindBy(ctx context.Context, phone string) (entities.User, error) {
	var u entities.User
	err := a.db.First(&u, "phone_number = ? and type= ?", phone, "admin").Error
	return u, err
}

func (a AuthenticateRepository) FindByUserID(ctx context.Context, userID uint) (entities.User, error) {
	var u entities.User
	err := a.db.First(&u, "id = ? and type=?", userID, "admin").Error
	return u, err
}
