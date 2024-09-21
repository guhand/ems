package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/request"
	"ems/domain"
	"ems/utils"
	"fmt"
)

type leaveService struct {
	leaveRepository      domain.LeaveRepository
	departmentRepository domain.DepartmentRepository
	userRepository       domain.UserRepository
}

func NewLeaveService(leaveRepository domain.LeaveRepository, departmentRepository domain.DepartmentRepository, userRepository domain.UserRepository) domain.LeaveService {
	return &leaveService{leaveRepository, departmentRepository, userRepository}
}

func (s *leaveService) RequestLeave(departmentMemberID uint, req *request.RequestLeave) error {
	isDepartmentMemberExists, err := s.departmentRepository.IsDepartmentMemberExists(departmentMemberID)

	if err != nil {
		return err
	}

	if !isDepartmentMemberExists {
		return apperror.DataNotFoundError("Employee")
	}

	isLeaveExistsWithoutApproval, err := s.leaveRepository.IsLeaveExistsWithoutApproval(departmentMemberID)

	if err != nil {
		return err
	}

	if isLeaveExistsWithoutApproval {
		return fmt.Errorf(`Last leave request in the pending state. Kindly contact the TL`)
	}

	if err := s.leaveRepository.RequestLeave(departmentMemberID, req); err != nil {
		return err
	}

	return nil
}

func (s *leaveService) FetchOwnLeaves(departmentMemberID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	data, err := s.leaveRepository.FetchOwnLeaves(departmentMemberID, filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *leaveService) FetchDepartmentMemberLeaves(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	data, err := s.leaveRepository.FetchDepartmentMemberLeaves(departmentID, filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *leaveService) UpdateLeaveStatus(leaveID uint, approvedBy uint, req *request.UpdateLeaveStatus) error {
	if err := s.leaveRepository.UpdateLeaveStatus(leaveID, approvedBy, req); err != nil {
		return err
	}

	return nil
}

func (s *leaveService) FetchLeaves(filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	data, err := s.leaveRepository.FetchLeaves(filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *leaveService) UpdateLeaveRequest(leaveID uint, req *request.RequestLeave) error {
	isPermissionExists, err := s.leaveRepository.IsLeaveExistsWithID(leaveID)

	if err != nil {
		return err
	}

	if !isPermissionExists {
		return apperror.DataNotFoundError("Leave")
	}

	if err := s.leaveRepository.UpdateLeaveRequest(leaveID, req); err != nil {
		return err
	}

	return nil
}

func (s *leaveService) RemoveLeaveRequest(leaveID uint) error {
	isLeaveExists, err := s.leaveRepository.IsLeaveExistsWithID(leaveID)

	if err != nil {
		return err
	}

	if !isLeaveExists {
		return apperror.DataNotFoundError("Leave")
	}

	isLeaveExistsWithApproval, err := s.leaveRepository.IsLeaveExistsWithApproval(leaveID)

	if err != nil {
		return err
	}

	if isLeaveExistsWithApproval {
		return fmt.Errorf("Approved cannot be removed")
	}

	if err := s.leaveRepository.RemoveLeaveRequest(leaveID); err != nil {
		return err
	}

	return nil
}
