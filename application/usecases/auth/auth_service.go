package auth

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	AdminUserResponse "shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/domain/repositories"
	"shop/interfaces/http/requests/admin"
)

type AuthenticateService struct {
	authRepo repositories.AuthenticateRepositoryInterface
}

func NewAuthenticateService(authRepo repositories.AuthenticateRepositoryInterface) AuthenticateServiceInterface {
	return &AuthenticateService{
		authRepo: authRepo,
	}
}

func (a *AuthenticateService) Login(ctx context.Context, req *requests.LoginRequest) (*AdminUserResponse.User, domain_err.CustomError) {

	user, err := a.authRepo.FindBy(ctx, req.Mobile)
	if user.ID == 0 {
		log.Println("[findUserBy error ] : ", err)
		if err != nil {
			return &AdminUserResponse.User{}, domain_err.HandleError(err, domain_err.RecordNotFound)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &AdminUserResponse.User{}, domain_err.New(err.Error(), domain_err.MobileOrPasswordIsWrong, 404)
	}

	return AdminUserResponse.ToUserResponse(user), domain_err.CustomError{}
}
