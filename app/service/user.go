package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/domain"
	"ems/utils"
	"fmt"
)

type userService struct {
	userRepository       domain.UserRepository
	departmentRepository domain.DepartmentRepository
}

func NewUserService(userRepository domain.UserRepository,
	departmentRepository domain.DepartmentRepository) domain.UserService {

	return &userService{userRepository, departmentRepository}
}

func (s *userService) CreateUser(req *request.CreateUser) error {
	isUserCodeExists, err := s.userRepository.IsUserCodeExists(req.Code)

	if err != nil {
		return err
	}

	if isUserCodeExists {
		return apperror.UniqueKeyError("User code")
	}

	isEmailExists, err := s.userRepository.IsEmailExists(req.Email)

	if err != nil {
		return err
	}

	if isEmailExists {
		return apperror.UniqueKeyError("Email")
	}

	isMobileExists, err := s.userRepository.IsMobileNumberExists(req.Mobile)

	if err != nil {
		return err
	}

	if isMobileExists {
		return apperror.UniqueKeyError("Mobile")
	}

	hashedPassword, err := utils.HashPassword("123")

	if err != nil {
		return err
	}

	// isAadharNumberExists, err := s.userRepository.IsUserAadharNumberExists(req.AadharNumber)

	// if err != nil {
	// 	return err
	// }

	// if isAadharNumberExists {
	// 	return apperror.UniqueKeyError("Aadhar Number")
	// }

	// isPanNumberExists, err := s.userRepository.IsUserPanNumberExists(req.PanNumber)

	// if err != nil {
	// 	return err
	// }

	// if isPanNumberExists {
	// 	return apperror.UniqueKeyError("Pan Number")
	// }

	// isBankAccountNumberExists, err := s.userRepository.IsUserbankAccountNumberExists(req.BankAccountNumber)

	// if err != nil {
	// 	return err
	// }

	// if isBankAccountNumberExists {
	// 	return apperror.UniqueKeyError("Bank Account Number")
	// }

	if err := s.userRepository.CreateUser(req, hashedPassword); err != nil {
		return nil
	}

	return nil
}

func (s *userService) FetchUsers(filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	data, err := s.userRepository.FetchUsers(filters)

	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *userService) UpdateUser(userID uint, req *request.UpdateUser) error {
	isUserExists, err := s.userRepository.IsUserExists(userID)

	if err != nil {
		return err
	}

	if !isUserExists {
		return apperror.DataNotFoundError("User")
	}

	isUserCodeExists, err := s.userRepository.IsUserCodeExistsExceptID(userID, req.Code)

	if err != nil {
		return err
	}

	if isUserCodeExists {
		return apperror.UniqueKeyError("User code")
	}

	isEmailExists, err := s.userRepository.IsEmailExistsExceptID(userID, req.Email)

	if err != nil {
		return err
	}

	if isEmailExists {
		return apperror.UniqueKeyError("Email")
	}

	isMobileExists, err := s.userRepository.IsMobileNumberExistsExceptID(userID, req.Mobile)

	if err != nil {
		return err
	}

	if isMobileExists {
		return apperror.UniqueKeyError("Mobile")
	}

	if err := s.userRepository.UpdateUser(userID, req); err != nil {
		return err
	}

	return nil
}

func (s *userService) FetchUserStatus() ([]response.FetchUserStatus, error) {
	data, err := s.userRepository.FetchUserStatus()

	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *userService) FetchLastUserCode() (*response.FetchLastUserCode, error) {
	data, err := s.userRepository.FetchLastUserCode()

	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *userService) FetchUnmappedLeadUsers() ([]response.FetchUnmappedLeadUsers, error) {
	data, err := s.userRepository.FetchUnmappedLeadUsers()

	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *userService) FetchUnmappedLeadUserIncludeUserID(req *request.FetchUnmappedLeadUserIncludeUserID) ([]response.FetchUnmappedLeadUsers, error) {
	if req.DepartmentID == 1 { // HR Department
		data, err := s.userRepository.FetchUnmappedHRUserIncludeUserID(req)

		if err != nil {
			return nil, err
		}

		return data, err
	}

	data, err := s.userRepository.FetchUnmappedLeadUserIncludeUserID(req)

	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *userService) FetchUnmappedUsers(req *request.FetchUnmappedUsersByDepartmentID) ([]response.FetchUnmappedEmployeeUsers, error) {
	isDepartmentExists, err := s.departmentRepository.IsDepartmentExists(req.DepartmentID)

	if err != nil {
		return nil, err
	}

	if !isDepartmentExists {
		return nil, apperror.DataNotFoundError("Department")
	}

	if req.DepartmentID == 1 { // HR Department
		data, err := s.userRepository.FetchUnmappedHRUsers()

		if err != nil {
			return nil, err
		}

		return data, err

	} else {
		data, err := s.userRepository.FetchUnmappedEmployeeUsers()

		if err != nil {
			return nil, err
		}

		return data, err
	}
}

func (s *userService) RemoveUser(userID uint) error {
	isMappedLeadUser, err := s.userRepository.IsMappedLeadUser(userID)

	if err != nil {
		return err
	}

	if isMappedLeadUser {
		return fmt.Errorf(`department lead mapped to department`)
	}

	if err := s.userRepository.RemoveUser(userID); err != nil {
		return err
	}

	return nil
}

func (s *userService) ResetPassword(userID uint, req *request.ResetPassword) error {
	isUserExists, err := s.userRepository.IsUserExists(userID)

	if err != nil {
		return err
	}

	if !isUserExists {
		return apperror.DataNotFoundError("User")
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return err
	}

	if err := s.userRepository.UpdatePassword(userID, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (s *userService) UpdateUserDetails(req *request.UpdateUserDetails) error {
	isUserExists, err := s.userRepository.IsUserExists(req.UserID)

	if err != nil {
		return err
	}

	if !isUserExists {
		return apperror.DataNotFoundError("User")
	}

	isAadharNumberExists, err := s.userRepository.IsAadharNumberExistsExceptID(req.UserID, req.AadharNumber)

	if err != nil {
		return err
	}

	if isAadharNumberExists {
		return apperror.UniqueKeyError("Aadhar Number")
	}

	isPanNumberExists, err := s.userRepository.IsPanNumberExistsExceptID(req.UserID, req.PanNumber)

	if err != nil {
		return err
	}

	if isPanNumberExists {
		return apperror.UniqueKeyError("Pan Number")
	}

	if err := s.userRepository.UpdateUserDetails(req); err != nil {
		return err
	}

	return nil
}

func (s *userService) FetchUserDetails(req *request.FetchUserDetails) (*response.FetchUserDetails, error) {
	isUserExists, err := s.userRepository.IsUserExists(req.UserID)

	if err != nil {
		return nil, err
	}

	if !isUserExists {
		return nil, apperror.DataNotFoundError("User")
	}

	data, err := s.userRepository.FetchUserDetails(req)

	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *userService) UploadFiles(userID uint, filepaths []string) error {
	isUserExists, err := s.userRepository.IsUserExists(userID)

	if err != nil {
		return err
	}

	if !isUserExists {
		return apperror.DataNotFoundError("User")
	}

	if err := s.userRepository.UploadFiles(userID, filepaths); err != nil {
		return err
	}

	return nil
}

func (s *userService) FetchFilePathsByUserID(userID uint) ([]response.FetchUploadedDocumentPaths, error) {
	isUserExists, err := s.userRepository.IsUserExists(userID)

	if err != nil {
		return nil, err
	}

	if !isUserExists {
		return nil, apperror.DataNotFoundError("User")
	}

	data, err := s.userRepository.FetchFilePathsByUserID(userID)

	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *userService) RemoveFile(filePath string) error {
	if err := s.userRepository.RemoveFile(filePath); err != nil {
		return err
	}
	return nil
}
