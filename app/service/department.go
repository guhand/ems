package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/constant"
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
		return apperror.UniqueKeyError("department name")
	}

	isDepartLeadUser, err := s.userRepository.IsUnmappedLeadUser(req.LeadID)

	if err != nil {
		return err
	}

	if !isDepartLeadUser {
		return apperror.DataNotFoundError("department lead user")
	}

	if err := s.departmentRepository.CreateDepartment(req); err != nil {
		return err
	}

	return nil
}

func (s *departmentService) FetchDepartments(filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	data, err := s.departmentRepository.FetchDepartments(filters)

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
		return apperror.DataNotFoundError("department")
	}

	isDepartmentNameExists, err := s.departmentRepository.IsDepartmentNameExistsExceptID(departmentID, req.Name)

	if err != nil {
		return err
	}

	if isDepartmentNameExists {
		return apperror.UniqueKeyError("department name")
	}

	if departmentID == 1 { //HR Department
		isUnmappedHRUserIncludeUserID, err := s.userRepository.IsUnmappedHRUserIncludeUserID(req.LeadID)

		if err != nil {
			return err
		}

		if !isUnmappedHRUserIncludeUserID {
			return apperror.DataNotFoundError("hr user")
		}
	} else {
		isUnmappedLeadUserIncludeUserID, err := s.userRepository.IsUnmappedLeadUserIncludeUserID(req.LeadID)

		if err != nil {
			return err
		}

		if !isUnmappedLeadUserIncludeUserID {
			return apperror.DataNotFoundError("department lead user")
		}
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
		return apperror.DataNotFoundError("department")
	}

	if departmentID == 1 { // HR Department
		return fmt.Errorf("hr department cannot be removed")
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
		return apperror.DataNotFoundError("department")
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
		return nil, apperror.DataNotFoundError("department")
	}

	data, err := s.departmentRepository.FetchDepartmentMembers(departmentID, filters)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *departmentService) UnMapUser(req *request.UnMapUser) error {
	isDepartmentUserExists, err := s.userRepository.IsDepartmentUserExists(req.UserID)

	if err != nil {
		return err
	}

	if isDepartmentUserExists.Count == 0 {
		return apperror.DataNotFoundError("department user")
	}

	if isDepartmentUserExists.RoleID == int(constant.DepartmentLead) {
		if req.LeadID == nil {
			return fmt.Errorf("please provide another lead to unmap this lead")
		} else {
			isUnmappedLeadUser, err := s.userRepository.IsUnmappedLeadUser(*req.LeadID)

			if err != nil {
				return err
			}

			if !isUnmappedLeadUser {
				return apperror.DataNotFoundError("department lead user")
			}

			if err := s.departmentRepository.UnMapUser(req); err != nil {
				return err
			}

			if err := s.departmentRepository.MapLeadToDepartment(uint(isDepartmentUserExists.DepartmentID),
				*req.LeadID); err != nil {
				return err
			}
		}
	} else if isDepartmentUserExists.RoleID == int(constant.HR) {
		departmentMemberCount, err := s.departmentRepository.GetDepartmentMemberCount(uint(isDepartmentUserExists.DepartmentID))
		if err != nil {
			return err
		}

		if departmentMemberCount == 1 && req.LeadID == nil {
			return fmt.Errorf("please provide another hr to unmap this hr")
		}

		if err := s.departmentRepository.UnMapUser(req); err != nil {
			return err
		}

		if departmentMemberCount == 1 && req.LeadID != nil {
			if err := s.departmentRepository.MapLeadToDepartment(uint(isDepartmentUserExists.DepartmentID),
				*req.LeadID); err != nil {
				return err
			}
		}
	} else {
		if err := s.departmentRepository.UnMapUser(req); err != nil {
			return err
		}
	}
	return nil
}
