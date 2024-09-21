package domain

import (
	"ems/app/model/request"
	"ems/utils"
)

type PermissionService interface {
	RequestPermission(departmentMemberID uint, req *request.RequestPermission) error
	FetchPermissions(departmentMemberID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	FetchUserPermissions(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdatePermissionStatus(permissionID, approvedBy uint, req *request.UpdatePermissionStatus) error
	UpdatePermissionRequest(permissionID uint, req *request.RequestPermission) error
	RemovePermissionRequest(permissionID uint) error
}

type PermissionRepository interface {
	RequestPermission(departmentMemberID uint, req *request.RequestPermission) error
	FetchPermissions(departmentMemberID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	FetchUserPermissions(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdatePermissionStatus(permissionID, approvedBy uint, req *request.UpdatePermissionStatus) error
	IsPermissionExistWithoutApproval(departmentMemberID uint) (bool, error)
	IsPermissionExistWithID(id uint) (bool, error)
	UpdatePermissionRequest(permissionID uint, req *request.RequestPermission) error
	RemovePermissionRequest(permissionID uint) error
	IsPermissionExistsWithApproval(permissionID uint) (bool, error)
	GetPermissionCount() (int, error)
	GetPermissionCountByDepartment(departmentID uint) (int, error)
	GetPermissionCountByUser(departmentMemberID uint) (int, error)
	GetApprovedPermissionCount() (int, error)
	GetApprovedPermissionCountByDepartment(departmentID uint) (int, error)
	GetApprovedPermissionCountByUser(departmentMemberID uint) (int, error)
}
