package auth

import (
	"context"
	AdminUserResponse "shop/application/dto/admin"
	"shop/domain/domain_err"
	"shop/interfaces/http/requests/admin"
)

type AuthenticateServiceInterface interface {
	Login(ctx context.Context, req *requests.LoginRequest) (*AdminUserResponse.User, domain_err.CustomError)
}
