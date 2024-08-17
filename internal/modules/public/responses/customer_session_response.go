package responses

import (
	"shop/internal/entities"
	"time"
)

type CustomerSession struct {
	ID         uint
	Mobile     string
	CustomerID uint
	SessionID  string
	IsActive   bool
	ExpiredAt  time.Time
}

func ToCustomerSession(session entities.Session) CustomerSession {
	return CustomerSession{
		ID:         session.ID,
		Mobile:     session.Mobile,
		CustomerID: session.CustomerID,
		SessionID:  session.SessionID,
		IsActive:   session.IsActive,
		ExpiredAt:  session.ExpiredAt,
	}
}
