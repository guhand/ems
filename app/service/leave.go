package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/constant"
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
		return apperror.DataNotFoundError("user")
	}

	for _, d := range req.Dates {
		_, isValidDate := utils.IsValidDate(d.Date)
		if !isValidDate {
			return fmt.Errorf("invalid date format: %s", d.Date)
		}
	}

	isLeaveExistsWithoutApproval, err := s.leaveRepository.IsLeaveExistsWithoutApproval(departmentMemberID)

	if err != nil {
		return err
	}

	if isLeaveExistsWithoutApproval && req.RoleID == uint(constant.Employee) {
		return fmt.Errorf("last leave request is in the pending state, please contact the TL")
	} else if isLeaveExistsWithoutApproval && (req.RoleID == uint(constant.DepartmentLead) || req.RoleID == uint(constant.HR)) {
		return fmt.Errorf(`last leave request is in the pending state. please contact the Manager`)
	}

	if err := s.leaveRepository.RequestLeave(departmentMemberID, req); err != nil {
		return err
	}

	return nil
}

func (s *leaveService) FetchOwnLeaves(departmentMemberID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	data, err := s.leaveRepository.FetchOwnLeaves(departmentMemberID, filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *leaveService) FetchDepartmentMemberLeaves(departmentID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
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

func (s *leaveService) FetchLeadAndHRLeaves(filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	data, err := s.leaveRepository.FetchLeadAndHRLeaves(filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *leaveService) UpdateLeaveRequest(leaveID uint, req *request.RequestLeave) error {
	isLeaveRequestExists, err := s.leaveRepository.IsLeaveExistsWithID(leaveID)

	if err != nil {
		return err
	}

	if !isLeaveRequestExists {
		return apperror.DataNotFoundError("leave request")
	}

	for _, d := range req.Dates {
		_, isValidDate := utils.IsValidDate(d.Date)
		if !isValidDate {
			return fmt.Errorf("invalid date format: %s", d.Date)
		}
	}

	if err := s.leaveRepository.UpdateLeaveRequest(leaveID, req); err != nil {
		return err
	}

	return nil
}

func (s *leaveService) RemoveLeaveRequest(leaveID uint) error {
	isLeaveRequestExists, err := s.leaveRepository.IsLeaveExistsWithID(leaveID)

	if err != nil {
		return err
	}

	if !isLeaveRequestExists {
		return apperror.DataNotFoundError("leave request")
	}

	isLeaveExistsWithApproval, err := s.leaveRepository.IsLeaveExistsWithApproval(leaveID)

	if err != nil {
		return err
	}

	if isLeaveExistsWithApproval {
		return fmt.Errorf("approved leave request cannot be removed")
	}

	if err := s.leaveRepository.RemoveLeaveRequest(leaveID); err != nil {
		return err
	}

	return nil
}
