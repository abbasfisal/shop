package entities

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Mobile     string
	CustomerID uint
	SessionID  string
	IsActive   bool
	ExpiredAt  time.Time

	//relation
	Customer Customer `gorm:"references:ID"`
}
