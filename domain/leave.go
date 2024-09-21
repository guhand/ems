package domain

import (
	"ems/app/model/request"
	"ems/utils"
)

type LeaveService interface {
	RequestLeave(departmentMemberID uint, req *request.RequestLeave) error
	FetchOwnLeaves(departmentMemberID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	FetchDepartmentMemberLeaves(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdateLeaveStatus(leaveID, approvedBy uint, req *request.UpdateLeaveStatus) error
	FetchLeaves(filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdateLeaveRequest(leaveID uint, req *request.RequestLeave) error
	RemoveLeaveRequest(leaveID uint) error
}

type LeaveRepository interface {
	RequestLeave(departmentMemberID uint, req *request.RequestLeave) error
	FetchOwnLeaves(departmentMemberID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	FetchDepartmentMemberLeaves(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdateLeaveStatus(leaveID, approvedBy uint, req *request.UpdateLeaveStatus) error
	IsLeaveExistsWithoutApproval(departmentMemberID uint) (bool, error)
	IsLeaveExistsWithApproval(leaveID uint) (bool, error)
	FetchLeaves(filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdateLeaveRequest(leaveID uint, req *request.RequestLeave) error
	IsLeaveExistsWithID(leaveID uint) (bool, error)
	RemoveLeaveRequest(leaveID uint) error
	GetLeaveCountByDepartment(departmentID uint) (float64, error)
	GetLeaveCountByUser(departmentMemberID uint) (float64, error)
	GetLeaveCount() (float64, error)
	GetApprovedLeaveCount() (float64, error)
	GetApprovedLeaveCountByUser(departmentMemberID uint) (float64, error)
	GetApprovedLeaveCountByDepartment(departmentID uint) (float64, error)
}
