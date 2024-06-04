package auth

import (
	"context"
	"gorm.io/gorm"
	"shop/internal/database/mysql"
	"shop/internal/entities"
)

type AuthenticateRepository struct {
	db *gorm.DB
}

func NewAuthenticateRepository() AuthenticateRepository {
	return AuthenticateRepository{
		db: mysql.Get(),
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
