package auth

import (
	"context"
	"shop/internal/modules/admin/requests"
	AdminUserResponse "shop/internal/modules/admin/responses"
	"shop/internal/pkg/custom_error"
)

type AuthenticateServiceInterface interface {
	Login(ctx context.Context, req requests.LoginRequest) (AdminUserResponse.User, custom_error.CustomError)
}
