package dashboard

import (
	"shop/domain/entities"
	"shop/domain/repositories"
)

type DashboardService struct {
	repo repositories.DashboardRepositoryInterface
}

func NewDashboardService(repo repositories.DashboardRepositoryInterface) *DashboardService {
	return &DashboardService{repo: repo}
}

func (d *DashboardService) GetStaticalData() (*entities.DashboardData, error) {
	return d.repo.GetDashboardStates()
}
