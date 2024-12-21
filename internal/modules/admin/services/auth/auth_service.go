package auth

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	"shop/internal/modules/admin/repositories/auth"
	"shop/internal/modules/admin/requests"
	AdminUserResponse "shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type AuthenticateService struct {
	authRepo auth.AuthenticateRepositoryInterface
}

func NewAuthenticateService(authRepo auth.AuthenticateRepositoryInterface) AuthenticateServiceInterface {
	return &AuthenticateService{
		authRepo: authRepo,
	}
}

func (a *AuthenticateService) Login(ctx context.Context, req *requests.LoginRequest) (*AdminUserResponse.User, custom_error.CustomError) {

	user, err := a.authRepo.FindBy(ctx, req.Mobile)
	if user.ID == 0 {
		log.Println("[findUserBy error ] : ", err)
		if err != nil {
			return &AdminUserResponse.User{}, custom_error.HandleError(err, custom_error.RecordNotFound)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &AdminUserResponse.User{}, custom_error.New(err.Error(), custom_error.MobileOrPasswordIsWrong, 404)
	}

	return AdminUserResponse.ToUserResponse(user), custom_error.CustomError{}
}
