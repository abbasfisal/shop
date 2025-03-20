package dashboard

import "shop/internal/modules/admin/repositories/dashboard"

type DashboardService struct {
	repo dashboard.DashboardRepositoryInterface
}

func NewDashboardService(repo dashboard.DashboardRepositoryInterface) *DashboardService {
	return &DashboardService{repo: repo}
}

func (d *DashboardService) GetStaticalData() (*dashboard.DashboardData, error) {
	return d.repo.GetDashboardStates()
}
