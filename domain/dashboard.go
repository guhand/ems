package domain

type DashboardService interface {
	FetchDashboardCounts(roleID, departmentMemberID, departmentID uint) (interface{}, error)
}
