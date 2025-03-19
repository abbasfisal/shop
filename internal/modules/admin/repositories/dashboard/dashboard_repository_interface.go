package dashboard

type DashboardRepositoryInterface interface {
	GetDashboardStates() (*DashboardData, error)
}
