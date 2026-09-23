package entities

import (
	"gorm.io/gorm"
	"time"
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
