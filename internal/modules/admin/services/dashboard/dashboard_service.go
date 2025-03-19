package main

import "shop/internal/modules/admin/repositories/dashboard"

type DashboardService struct {
	repo dashboard.DashboardRepositoryInterface
}

func NewDashboardService(repo dashboard.DashboardRepositoryInterface) *DashboardService {
	return &DashboardService{repo: repo}
}
