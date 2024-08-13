package entities

import "gorm.io/gorm"

type OTP struct {
	gorm.Model
	Mobile    string
	Code      string
	IsExpired bool
}
