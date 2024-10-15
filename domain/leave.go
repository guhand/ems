package domain

import (
	"ems/app/model/request"
	"ems/utils"
)

type LeaveService interface {
	RequestLeave(departmentMemberID uint, req *request.RequestLeave) error
	FetchOwnLeaves(departmentMemberID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	FetchDepartmentMemberLeaves(departmentID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	UpdateLeaveStatus(leaveID, approvedBy uint, req *request.UpdateLeaveStatus) error
	FetchLeadAndHRLeaves(filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	UpdateLeaveRequest(leaveID uint, req *request.RequestLeave) error
	RemoveLeaveRequest(leaveID uint) error
}

type LeaveRepository interface {
	RequestLeave(departmentMemberID uint, req *request.RequestLeave) error
	FetchOwnLeaves(departmentMemberID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	FetchDepartmentMemberLeaves(departmentID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	UpdateLeaveStatus(leaveID, approvedBy uint, req *request.UpdateLeaveStatus) error
	IsLeaveExistsWithoutApproval(departmentMemberID uint) (bool, error)
	IsLeaveExistsWithApproval(leaveID uint) (bool, error)
	FetchLeadAndHRLeaves(filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	UpdateLeaveRequest(leaveID uint, req *request.RequestLeave) error
	IsLeaveExistsWithID(leaveID uint) (bool, error)
	RemoveLeaveRequest(leaveID uint) error
	GetLeaveCountByUser(departmentMemberID uint, dateFilters *request.DateFilters) (float64, error)
	GetLeaveCount(dateFilters *request.DateFilters) (float64, error)
	GetApprovedLeaveCount(dateFilters *request.DateFilters) (float64, error)
	GetApprovedLeaveCountByUser(departmentMemberID uint, dateFilters *request.DateFilters) (float64, error)
}
