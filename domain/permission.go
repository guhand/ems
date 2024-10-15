package domain

import (
	"ems/app/model/request"
	"ems/utils"
)

type PermissionService interface {
	RequestPermission(departmentMemberID uint, req *request.RequestPermission) error
	FetchOwnPermissions(departmentMemberID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	FetchDepartmentMemberPermissions(departmentID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	UpdatePermissionStatus(permissionID, approvedBy uint, req *request.UpdatePermissionStatus) error
	UpdatePermissionRequest(departmentMemberID uint, permissionID uint, req *request.RequestPermission) error
	RemovePermissionRequest(permissionID uint) error
	FetchLeadAndHRPermissions(filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
}

type PermissionRepository interface {
	RequestPermission(departmentMemberID uint, req *request.RequestPermission) error
	FetchOwnPermissions(departmentMemberID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	FetchDepartmentMemberPermissions(departmentID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
	UpdatePermissionStatus(permissionID, approvedBy uint, req *request.UpdatePermissionStatus) error
	IsPermissionExistWithoutApproval(departmentMemberID uint) (bool, error)
	IsPermissionExistWithID(id uint) (bool, error)
	UpdatePermissionRequest(permissionID uint, req *request.RequestPermission) error
	RemovePermissionRequest(permissionID uint) error
	IsPermissionExistsWithApproval(permissionID uint) (bool, error)
	GetPermissionCount(dateFilters *request.DateFilters) (int, error)
	GetPermissionCountByUser(departmentMemberID uint, dateFilters *request.DateFilters) (int, error)
	GetApprovedPermissionCount(dateFilters *request.DateFilters) (int, error)
	GetApprovedPermissionCountByUser(departmentMemberID uint, dateFilters *request.DateFilters) (int, error)
	FetchLeadAndHRPermissions(filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error)
}
