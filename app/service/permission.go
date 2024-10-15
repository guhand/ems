package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/constant"
	"ems/app/model/request"
	"ems/domain"
	"ems/utils"
	"fmt"
)

type permissionService struct {
	permissionRepository domain.PermissionRepository
	departmentRepository domain.DepartmentRepository
	userRepository       domain.UserRepository
}

func NewPermissionService(permissionRepository domain.PermissionRepository, departmentRepository domain.DepartmentRepository, userRepository domain.UserRepository) domain.PermissionService {
	return &permissionService{permissionRepository, departmentRepository, userRepository}
}

func (s *permissionService) RequestPermission(departmentMemberID uint, req *request.RequestPermission) error {
	var dateFilters request.DateFilters

	isDepartmentMemberExists, err := s.departmentRepository.IsDepartmentMemberExists(departmentMemberID)

	if err != nil {
		return err
	}

	if !isDepartmentMemberExists {
		return apperror.DataNotFoundError("user")
	}

	date, isValidDate := utils.IsValidDate(req.Date)

	if !isValidDate {
		return fmt.Errorf("invalid date format: %s", req.Date)
	}

	dateFilters.Year = date.Year()
	dateFilters.Month = int(date.Month())

	if date.Day() > 26 {
		dateFilters.Month = int(date.Month()) + 1
		if date.Month() == 12 {
			dateFilters.Month = 1
			dateFilters.Year = date.Year() + 1
		}
	}

	if err := utils.ValidateTimeDifference(req.FromTime, req.ToTime); err != nil {
		return err
	}

	permissionCount, err := s.permissionRepository.GetPermissionCountByUser(departmentMemberID, &dateFilters)

	if err != nil {
		return err
	}

	if permissionCount >= 3 {
		return fmt.Errorf("permission limit exceeded: maximum of 3 permissions reached for this month")
	}

	isPermissionExistsWithoutApproval, err := s.permissionRepository.IsPermissionExistWithoutApproval(departmentMemberID)

	if err != nil {
		return err
	}

	if isPermissionExistsWithoutApproval && req.RoleID == uint(constant.Employee) {
		return fmt.Errorf(`last permission request is in the pending state. please contact the TL`)
	} else if isPermissionExistsWithoutApproval && (req.RoleID == uint(constant.DepartmentLead) || req.RoleID == uint(constant.HR)) {
		return fmt.Errorf(`last permission request is in the pending state. please contact the Manager`)
	}

	if err := s.permissionRepository.RequestPermission(departmentMemberID, req); err != nil {
		return err
	}

	return nil
}

func (s *permissionService) FetchOwnPermissions(departmentMemberID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	data, err := s.permissionRepository.FetchOwnPermissions(departmentMemberID, filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *permissionService) FetchDepartmentMemberPermissions(departmentID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	data, err := s.permissionRepository.FetchDepartmentMemberPermissions(departmentID, filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *permissionService) UpdatePermissionStatus(permissionID uint, approvedBy uint, req *request.UpdatePermissionStatus) error {
	if err := s.permissionRepository.UpdatePermissionStatus(permissionID, approvedBy, req); err != nil {
		return err
	}

	return nil
}

func (s *permissionService) UpdatePermissionRequest(departmentMemberID uint, permissionID uint, req *request.RequestPermission) error {
	isPermissionExists, err := s.permissionRepository.IsPermissionExistWithID(permissionID)

	if err != nil {
		return err
	}

	if !isPermissionExists {
		return apperror.DataNotFoundError("permission")
	}

	date, isValidDate := utils.IsValidDate(req.Date)

	if !isValidDate {
		return fmt.Errorf("invalid date format: %s", req.Date)
	}

	var dateFilters request.DateFilters
	dateFilters.Month = int(date.Month())
	dateFilters.Year = date.Year()

	if date.Day() > 26 {
		dateFilters.Month = int(date.Month()) + 1
		if date.Month() == 12 {
			dateFilters.Month = 1
			dateFilters.Year = date.Year() + 1
		}
	}

	permissionCount, err := s.permissionRepository.GetPermissionCountByUser(departmentMemberID, &dateFilters)

	if err != nil {
		return err
	}

	if permissionCount >= 3 {
		return fmt.Errorf("permission limit exceeded: maximum of 3 permissions reached for this month")
	}

	if err := utils.ValidateTimeDifference(req.FromTime, req.ToTime); err != nil {
		return err
	}

	if err := utils.ValidateTimeDifference(req.FromTime, req.ToTime); err != nil {
		return err
	}

	if err := s.permissionRepository.UpdatePermissionRequest(permissionID, req); err != nil {
		return err
	}

	return nil
}

func (s *permissionService) RemovePermissionRequest(permissionID uint) error {
	isPermissionExists, err := s.permissionRepository.IsPermissionExistWithID(permissionID)

	if err != nil {
		return err
	}

	if !isPermissionExists {
		return apperror.DataNotFoundError("permission")
	}

	isPermissionExistsWithApproval, err := s.permissionRepository.IsPermissionExistsWithApproval(permissionID)

	if err != nil {
		return err
	}

	if isPermissionExistsWithApproval {
		return fmt.Errorf("approved permission request cannot be removed")
	}

	if err := s.permissionRepository.RemovePermissionRequest(permissionID); err != nil {
		return err
	}

	return nil
}

func (s *permissionService) FetchLeadAndHRPermissions(filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	data, err := s.permissionRepository.FetchLeadAndHRPermissions(filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}
