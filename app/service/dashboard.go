package service

import (
	"ems/app/model/constant"
	"ems/app/model/request"
	"ems/domain"
)

type dashboardService struct {
	userRepository       domain.UserRepository
	departRepository     domain.DepartmentRepository
	leaveRepository      domain.LeaveRepository
	permissionRepository domain.PermissionRepository
	noticeRepository     domain.NoticeRepository
}

func NewDashboardService(userRepository domain.UserRepository,
	departRepository domain.DepartmentRepository, leaveRepository domain.LeaveRepository,
	permissionRepository domain.PermissionRepository,
	noticeRepository domain.NoticeRepository) domain.DashboardService {

	return &dashboardService{userRepository, departRepository, leaveRepository,
		permissionRepository, noticeRepository}
}

func (s *dashboardService) FetchDashboardCounts(roleID, departmentMemberID, departmentID uint, dateFilters *request.DateFilters) (interface{}, error) {

	var data struct {
		EmployeeCount           *int     `json:"employeeCount,omitempty"`
		EmployeeOnNoticeCount   *int     `json:"employeeOnNoticeCount,omitempty"`
		DepartmentCount         *int     `json:"departmentCount,omitempty"`
		LeaveCount              *float64 `json:"leaveCount,omitempty"`
		ApprovedLeaveCount      *float64 `json:"approvedLeaveCount,omitempty"`
		PendingLeaveCount       *float64 `json:"pendingLeaveCount,omitempty"`
		PermissionCount         *int     `json:"permissionCount,omitempty"`
		ApprovedPermissionCount *int     `json:"approvedPermissionCount,omitempty"`
		PendingPermissionCount  *int     `json:"pendingPermissionCount,omitempty"`
	}

	switch int(roleID) {
	case int(constant.Admin), int(constant.Manager), int(constant.HR):
		{
			employeeCount, err := s.userRepository.GetUserCount()
			if err != nil {
				return nil, err
			}

			employeeOnNoticeCount, err := s.noticeRepository.GetNoticeUserCount()
			if err != nil {
				return nil, err
			}

			leaveCount, err := s.leaveRepository.GetLeaveCount(dateFilters)
			if err != nil {
				return nil, err
			}

			approvedLeaveCount, err := s.leaveRepository.GetApprovedLeaveCount(dateFilters)
			if err != nil {
				return nil, err
			}

			permissionCount, err := s.permissionRepository.GetPermissionCount(dateFilters)
			if err != nil {
				return nil, err
			}

			approvedPermissionCount, err := s.permissionRepository.GetApprovedPermissionCount(dateFilters)
			if err != nil {
				return nil, err
			}

			data.EmployeeCount = &employeeCount
			data.EmployeeOnNoticeCount = &employeeOnNoticeCount
			data.LeaveCount = &leaveCount
			data.ApprovedLeaveCount = &approvedLeaveCount

			pendingLeaveCount := leaveCount - approvedLeaveCount
			data.PendingLeaveCount = &pendingLeaveCount

			data.PermissionCount = &permissionCount
			data.ApprovedPermissionCount = &approvedPermissionCount

			pendingPermissionCount := permissionCount - approvedPermissionCount
			data.PendingPermissionCount = &pendingPermissionCount
		}

	case int(constant.DepartmentLead):
		{
			employeeCount, err := s.departRepository.GetDepartmentMemberCount(departmentID)
			if err != nil {
				return nil, err
			}

			employeeOnNoticeCount, err := s.noticeRepository.GetNoticeUserCountByDepartment(departmentID)
			if err != nil {
				return nil, err
			}

			leaveCount, err := s.leaveRepository.GetLeaveCountByUser(departmentMemberID, dateFilters)
			if err != nil {
				return nil, err
			}

			approvedLeaveCount, err := s.leaveRepository.GetApprovedLeaveCountByUser(departmentMemberID, dateFilters)
			if err != nil {
				return nil, err
			}

			permissionCount, err := s.permissionRepository.GetPermissionCountByUser(departmentMemberID, dateFilters)
			if err != nil {
				return nil, err
			}

			approvedPermissionCount, err := s.permissionRepository.GetApprovedPermissionCountByUser(departmentMemberID, dateFilters)
			if err != nil {
				return nil, err
			}

			data.EmployeeCount = &employeeCount
			data.EmployeeOnNoticeCount = &employeeOnNoticeCount
			data.LeaveCount = &leaveCount
			data.ApprovedLeaveCount = &approvedLeaveCount

			pendingLeaveCount := leaveCount - approvedLeaveCount
			data.PendingLeaveCount = &pendingLeaveCount

			data.PermissionCount = &permissionCount
			data.ApprovedPermissionCount = &approvedPermissionCount

			pendingPermissionCount := permissionCount - approvedPermissionCount
			data.PendingPermissionCount = &pendingPermissionCount
		}

	case int(constant.Employee):
		{
			leaveCount, err := s.leaveRepository.GetLeaveCountByUser(departmentMemberID, dateFilters)
			if err != nil {
				return nil, err
			}

			approvedLeaveCount, err := s.leaveRepository.GetApprovedLeaveCountByUser(departmentMemberID, dateFilters)
			if err != nil {
				return nil, err
			}

			permissionCount, err := s.permissionRepository.GetPermissionCountByUser(departmentMemberID, dateFilters)
			if err != nil {
				return nil, err
			}

			approvedPermissionCount, err := s.permissionRepository.GetApprovedPermissionCountByUser(departmentMemberID, dateFilters)
			if err != nil {
				return nil, err
			}

			data.LeaveCount = &leaveCount
			data.ApprovedLeaveCount = &approvedLeaveCount

			pendingLeaveCount := leaveCount - approvedLeaveCount
			data.PendingLeaveCount = &pendingLeaveCount

			data.PermissionCount = &permissionCount
			data.ApprovedPermissionCount = &approvedPermissionCount

			pendingPermissionCount := permissionCount - approvedPermissionCount
			data.PendingPermissionCount = &pendingPermissionCount
		}
	}

	return data, nil
}
