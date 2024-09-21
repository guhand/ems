package domain

import (
	"ems/app/model/request"
	"ems/utils"
)

type DepartmentService interface {
	CreateDepartment(req *request.CreateDepartment) error
	FectchDepartments(filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdateDepartment(departmentID uint, req *request.UpdateDepartment) error
	RemoveDepartment(departmentID uint) error
	MappUsersToDepartment(departmentID uint, req *request.MappUsersToDepartment) error
	FetchDepartmentMembers(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UnMapEmployee(req *request.UnMapEmployee) error
}

type DepartmentRepository interface {
	CreateDepartment(req *request.CreateDepartment) error
	IsDepartmentNameExists(name string) (bool, error)
	IsDepartmentExists(id uint) (bool, error)
	FectchDepartments(filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdateDepartment(id uint, req *request.UpdateDepartment) error
	RemoveDepartment(departmentID uint) error
	IsDepartmentNameExistsExceptID(id uint, name string) (bool, error)
	MappUsersToDepartment(departmentID uint, req *request.MappUsersToDepartment) error
	FetchDepartmentMembers(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UnMapEmployee(req *request.UnMapEmployee) error
	IsDepartmentMemberExists(id uint) (bool, error)
	GetDepartmentCount() (int, error)
	GetDepartmentMemberCount(departmentID uint) (int, error)
}
