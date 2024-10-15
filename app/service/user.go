package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/domain"
	"ems/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository       domain.UserRepository
	departmentRepository domain.DepartmentRepository
	leaveRepository      domain.LeaveRepository
	permissionRepository domain.PermissionRepository
}

func NewUserService(userRepository domain.UserRepository,
	departmentRepository domain.DepartmentRepository, leaveRepository domain.LeaveRepository, permissionRepository domain.PermissionRepository) domain.UserService {

	return &userService{userRepository, departmentRepository, leaveRepository, permissionRepository}
}

func (s *userService) CreateUser(req *request.CreateUser) error {
	isUserCodeExists, err := s.userRepository.IsUserCodeExists(req.Code)

	if err != nil {
		return err
	}

	if isUserCodeExists {
		return apperror.UniqueKeyError("user code")
	}

	isEmailExists, err := s.userRepository.IsEmailExists(req.Email)

	if err != nil {
		return err
	}

	if isEmailExists {
		return apperror.UniqueKeyError("email")
	}

	isMobileExists, err := s.userRepository.IsMobileNumberExists(req.Mobile)

	if err != nil {
		return err
	}

	if isMobileExists {
		return apperror.UniqueKeyError("mobile")
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return err
	}

	if err := s.userRepository.CreateUser(req, hashedPassword); err != nil {
		return nil
	}

	return nil
}

func (s *userService) FetchUsers(filters *request.FetchUsers) (*utils.PaginationResponse, error) {
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
		return apperror.DataNotFoundError("user")
	}

	isUserCodeExists, err := s.userRepository.IsUserCodeExistsExceptID(userID, req.Code)

	if err != nil {
		return err
	}

	if isUserCodeExists {
		return apperror.UniqueKeyError("user code")
	}

	isEmailExists, err := s.userRepository.IsEmailExistsExceptID(userID, req.Email)

	if err != nil {
		return err
	}

	if isEmailExists {
		return apperror.UniqueKeyError("email")
	}

	isMobileExists, err := s.userRepository.IsMobileNumberExistsExceptID(userID, req.Mobile)

	if err != nil {
		return err
	}

	if isMobileExists {
		return apperror.UniqueKeyError("mobile")
	}

	if err := s.userRepository.UpdateUser(userID, req); err != nil {
		return err
	}

	return nil
}

func (s *userService) FetchLastUserCode() (*response.FetchLastUserCode, error) {
	data, err := s.userRepository.FetchLastUserCode()

	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *userService) FetchUnmappedLeadUsers() ([]response.FetchUnmappedUsers, error) {
	data, err := s.userRepository.FetchUnmappedLeadUsers()

	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *userService) FetchUnmappedLeadUserIncludeUserID(req *request.FetchUnmappedLeadUserIncludeUserID) ([]response.FetchUnmappedUsers, error) {
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

func (s *userService) FetchUnmappedUsers(req *request.FetchUnmappedUsersByDepartmentID) ([]response.FetchUnmappedUsers, error) {
	isDepartmentExists, err := s.departmentRepository.IsDepartmentExists(req.DepartmentID)

	if err != nil {
		return nil, err
	}

	if !isDepartmentExists {
		return nil, apperror.DataNotFoundError("department")
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
		return fmt.Errorf(`user mapped to department. Kindly unmap from department`)
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
		return apperror.DataNotFoundError("user")
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
		return apperror.DataNotFoundError("user")
	}

	isAadharNumberExists, err := s.userRepository.IsAadharNumberExistsExceptID(req.UserID, req.AadharNumber)

	if err != nil {
		return err
	}

	if isAadharNumberExists {
		return apperror.UniqueKeyError("aadhar Number")
	}

	isPanNumberExists, err := s.userRepository.IsPanNumberExistsExceptID(req.UserID, req.PanNumber)

	if err != nil {
		return err
	}

	if isPanNumberExists {
		return apperror.UniqueKeyError("pan Number")
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
		return nil, apperror.DataNotFoundError("user")
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
		return apperror.DataNotFoundError("user")
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
		return nil, apperror.DataNotFoundError("user")
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

func (s *userService) ChangePassword(req *request.ChangePassword) error {
	user, err := s.userRepository.GetUserByEmail(req.Email)

	if err != nil {
		return err
	}

	if user == nil {
		return apperror.DataNotFoundError("email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("incorrect old password")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)

	if err != nil {
		return err
	}

	if err := s.userRepository.UpdatePassword(user.ID, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (s *userService) FetchUnmappedHRUsers() ([]response.FetchUnmappedUsers, error) {

	data, err := s.userRepository.FetchUnmappedHRUsers()

	if err != nil {
		return nil, err
	}

	return data, err
}
