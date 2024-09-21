package service

import (
	apperror "ems/app/model/app_error"
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
	isDepartmentMemberExists, err := s.departmentRepository.IsDepartmentMemberExists(departmentMemberID)

	if err != nil {
		return err
	}

	if !isDepartmentMemberExists {
		return apperror.DataNotFoundError("Employee")
	}

	if err := utils.ValidateTimeDifference(req.FromTime, req.ToTime); err != nil {
		return err
	}

	permissionCount, err := s.permissionRepository.GetPermissionCountByUser(departmentMemberID)

	if err != nil {
		return err
	}

	if permissionCount >= 3 {
		return fmt.Errorf("Permission limit exceeded: You have reached the maximum of 3 permissions for this month")
	}

	isLeaveExistsWithoutApproval, err := s.permissionRepository.IsPermissionExistWithoutApproval(departmentMemberID)

	if err != nil {
		return err
	}

	if isLeaveExistsWithoutApproval {
		return fmt.Errorf(`Last permission request in the pending state. Kindly contact the TL`)
	}

	if err := s.permissionRepository.RequestPermission(departmentMemberID, req); err != nil {
		return err
	}

	return nil
}

func (s *permissionService) FetchPermissions(departmentMemberID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	data, err := s.permissionRepository.FetchPermissions(departmentMemberID, filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *permissionService) FetchUserPermissions(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	data, err := s.permissionRepository.FetchUserPermissions(departmentID, filters)

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

func (s *permissionService) UpdatePermissionRequest(permissionID uint, req *request.RequestPermission) error {
	isPermissionExists, err := s.permissionRepository.IsPermissionExistWithID(permissionID)

	if err != nil {
		return err
	}

	if !isPermissionExists {
		return apperror.DataNotFoundError("Permission")
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
		return apperror.DataNotFoundError("Permission")
	}

	isPermissionExistsWithApproval, err := s.permissionRepository.IsPermissionExistsWithApproval(permissionID)

	if err != nil {
		return err
	}

	if isPermissionExistsWithApproval {
		return fmt.Errorf("Approved permission cannot be removed")
	}

	if err := s.permissionRepository.RemovePermissionRequest(permissionID); err != nil {
		return err
	}

	return nil
}
