package auth

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"shop/internal/modules/admin/repositories/auth"
	"shop/internal/modules/admin/requests"
	AdminUserResponse "shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type AuthenticateService struct {
	authRepo auth.AuthenticateRepositoryInterface
}

func NewAuthenticateService() AuthenticateService {
	return AuthenticateService{
		authRepo: auth.NewAuthenticateRepository(),
	}
}

func (a AuthenticateService) Login(ctx context.Context, req requests.LoginRequest) (AdminUserResponse.User, custom_error.CustomError) {
	var customErr custom_error.CustomError
	var userResponse AdminUserResponse.User
	user, err := a.authRepo.FindBy(ctx, req.Mobile)

	if user.ID == 0 {
		log.Println("[findUserBy error ] : ", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userResponse, custom_error.New(err.Error(), custom_error.MobileOrPasswordIsWrong, 404)
		}
		return userResponse, custom_error.New(err.Error(), custom_error.InternalServerError, 500)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return userResponse, custom_error.New(err.Error(), custom_error.MobileOrPasswordIsWrong, 404)
	}

	return AdminUserResponse.ToUserResponse(user), customErr
}
