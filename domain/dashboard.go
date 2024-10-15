package domain

import "ems/app/model/request"

type DashboardService interface {
	FetchDashboardCounts(roleID, departmentMemberID, departmentID uint, dateFilters *request.DateFilters) (interface{}, error)
}
