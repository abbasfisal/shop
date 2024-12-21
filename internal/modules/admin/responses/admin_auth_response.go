package responses

import "shop/internal/entities"

type User struct {
	ID          uint
	FirstName   string
	LastName    string
	PhoneNumber string
	Type        string
}

func ToUserResponse(user *entities.User) *User {
	return &User{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Type:        user.Type,
	}
}
