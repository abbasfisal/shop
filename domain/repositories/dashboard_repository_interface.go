package repositories

import "shop/domain/entities"

type DashboardRepositoryInterface interface {
	GetDashboardStates() (*entities.DashboardData, error)
}
