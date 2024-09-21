package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/request"
	"ems/domain"
	"ems/utils"
	"fmt"
)

type departmentService struct {
	departmentRepository domain.DepartmentRepository
	userRepository       domain.UserRepository
}

func NewDepartmentService(departmentRepository domain.DepartmentRepository, userRepository domain.UserRepository) domain.DepartmentService {
	return &departmentService{departmentRepository, userRepository}
}

func (s *departmentService) CreateDepartment(req *request.CreateDepartment) error {

	isDepartmentNameExists, err := s.departmentRepository.IsDepartmentNameExists(req.Name)

	if err != nil {
		return err
	}

	if isDepartmentNameExists {
		return apperror.UniqueKeyError("Department name")
	}

	isDepartLeadUser, err := s.userRepository.IsUnmappedLeadUser(req.LeadID)

	if err != nil {
		return err
	}

	if !isDepartLeadUser {
		return apperror.DataNotFoundError("Department lead user")
	}

	if err := s.departmentRepository.CreateDepartment(req); err != nil {
		return err
	}

	return nil
}

func (s *departmentService) FectchDepartments(filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	data, err := s.departmentRepository.FectchDepartments(filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *departmentService) UpdateDepartment(departmentID uint, req *request.UpdateDepartment) error {
	isDepartmentExists, err := s.departmentRepository.IsDepartmentExists(departmentID)

	if err != nil {
		return err
	}

	if !isDepartmentExists {
		return apperror.DataNotFoundError("Department")
	}

	isDepartmentNameExists, err := s.departmentRepository.IsDepartmentNameExistsExceptID(departmentID, req.Name)

	if err != nil {
		return err
	}

	if isDepartmentNameExists {
		return apperror.UniqueKeyError("Department name")
	}

	isUnmappedLeadUserIncludeUserID, err := s.userRepository.IsUnmappedLeadUserIncludeUserID(req.LeadID)

	if err != nil {
		return err
	}

	if !isUnmappedLeadUserIncludeUserID {
		return apperror.DataNotFoundError("Department lead user")
	}

	if err := s.departmentRepository.UpdateDepartment(departmentID, req); err != nil {
		return err
	}

	return nil
}

func (s *departmentService) RemoveDepartment(departmentID uint) error {
	isDepartmentExists, err := s.departmentRepository.IsDepartmentExists(departmentID)

	if err != nil {
		return err
	}

	if !isDepartmentExists {
		return apperror.DataNotFoundError("Department")
	}

	if departmentID == 1 { // HR Department
		return fmt.Errorf("HR Department cannot be removed")
	}

	if err := s.departmentRepository.RemoveDepartment(departmentID); err != nil {
		return err
	}

	return nil
}

func (s *departmentService) MappUsersToDepartment(departmentID uint, req *request.MappUsersToDepartment) error {
	isDepartmentExists, err := s.departmentRepository.IsDepartmentExists(departmentID)

	if err != nil {
		return err
	}

	if !isDepartmentExists {
		return apperror.DataNotFoundError("Department")
	}

	if departmentID == 1 { // HR Department
		employeesCount, err := s.userRepository.GetUnmappedHRsCount(req.UserIDs)

		if err != nil {
			return err
		}

		if employeesCount != len(req.UserIDs) {
			return apperror.DataNotFoundError("some of the HRs")
		}
	} else {
		employeesCount, err := s.userRepository.GetUnmappedEmployeesCount(req.UserIDs)

		if err != nil {
			return err
		}

		if employeesCount != len(req.UserIDs) {
			return apperror.DataNotFoundError("some of the employees")
		}
	}

	if err := s.departmentRepository.MappUsersToDepartment(departmentID, req); err != nil {
		return err
	}

	return nil
}

func (s *departmentService) FetchDepartmentMembers(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	isDepartmentExists, err := s.departmentRepository.IsDepartmentExists(departmentID)

	if err != nil {
		return nil, err
	}

	if !isDepartmentExists {
		return nil, apperror.DataNotFoundError("Department")
	}

	data, err := s.departmentRepository.FetchDepartmentMembers(departmentID, filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *departmentService) UnMapEmployee(req *request.UnMapEmployee) error {
	isEmployeeExists, err := s.userRepository.IsEmployeeUserExists(req.UserID)

	if err != nil {
		return err
	}

	if !isEmployeeExists {
		return apperror.DataNotFoundError("Employee")
	}

	return s.departmentRepository.UnMapEmployee(req)
}
