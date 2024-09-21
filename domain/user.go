package domain

import (
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/app/model/schema"
	"ems/utils"
)

type UserService interface {
	CreateUser(req *request.CreateUser) error
	FetchUsers(filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdateUser(userID uint, req *request.UpdateUser) error
	FetchUserStatus() ([]response.FetchUserStatus, error)
	FetchLastUserCode() (*response.FetchLastUserCode, error)
	FetchUnmappedLeadUsers() ([]response.FetchUnmappedLeadUsers, error)
	FetchUnmappedLeadUserIncludeUserID(req *request.FetchUnmappedLeadUserIncludeUserID) ([]response.FetchUnmappedLeadUsers, error)
	FetchUnmappedUsers(req *request.FetchUnmappedUsersByDepartmentID) ([]response.FetchUnmappedEmployeeUsers, error)
	RemoveUser(userID uint) error
	ResetPassword(userID uint, req *request.ResetPassword) error
	UpdateUserDetails(req *request.UpdateUserDetails) error
	FetchUserDetails(req *request.FetchUserDetails) (*response.FetchUserDetails, error)
	UploadFiles(userID uint, filepaths []string) error
	FetchFilePathsByUserID(userID uint) ([]response.FetchUploadedDocumentPaths, error)
	RemoveFile(filePath string) error
}

type UserRepository interface {
	GetUserByEmail(email string) (*response.FetchUserByEmail, error)
	CreateOTP(data *schema.ForgotPasswordOtp) (bool, error)
	GetOTPStatusByUserID(id uint) (*schema.ForgotPasswordOtp, error)
	UpdateToken(userID uint, token string) error
	UpdateOTPStatus(userID uint) error
	GetUserByID(id uint) (*response.FetchUserByID, error)
	IsEmailExists(email string) (bool, error)
	IsMobileNumberExists(mobile string) (bool, error)
	IsEmailExistsExceptID(id uint, email string) (bool, error)
	IsMobileNumberExistsExceptID(id uint, mobile string) (bool, error)
	IsAadharNumberExistsExceptID(id uint, aadharNumber string) (bool, error)
	IsPanNumberExistsExceptID(id uint, panNumber string) (bool, error)
	IsUserCodeExists(code string) (bool, error)
	IsUserCodeExistsExceptID(id uint, code string) (bool, error)
	CreateUser(req *request.CreateUser, hashedPassword string) error
	FetchUsers(filters *request.CommonRequest) (*utils.PaginationResponse, error)
	UpdateUser(userID uint, req *request.UpdateUser) error
	IsUserExists(userID uint) (bool, error)
	FetchUserStatus() ([]response.FetchUserStatus, error)
	IsUnmappedLeadUser(userID uint) (bool, error)
	IsUnmappedLeadUserIncludeUserID(userID uint) (bool, error)
	FetchLastUserCode() (*response.FetchLastUserCode, error)
	UpdateTokenStatus(userID uint) error
	UpdatePassword(userId uint, hashedPassword string) error
	FetchUnmappedLeadUsers() ([]response.FetchUnmappedLeadUsers, error)
	FetchUnmappedLeadUserIncludeUserID(req *request.FetchUnmappedLeadUserIncludeUserID) ([]response.FetchUnmappedLeadUsers, error)
	FetchUnmappedHRUserIncludeUserID(req *request.FetchUnmappedLeadUserIncludeUserID) ([]response.FetchUnmappedLeadUsers, error)
	GetUnmappedEmployeesCount(userIDs []uint) (int, error)
	GetUnmappedHRsCount(userIDs []uint) (int, error)
	FetchUnmappedEmployeeUsers() ([]response.FetchUnmappedEmployeeUsers, error)
	FetchUnmappedHRUsers() ([]response.FetchUnmappedEmployeeUsers, error)
	IsEmployeeUserExists(id uint) (bool, error)
	RemoveUser(userID uint) error
	IsMappedLeadUser(userID uint) (bool, error)
	UpdateUserDetails(req *request.UpdateUserDetails) error
	FetchUserDetails(req *request.FetchUserDetails) (*response.FetchUserDetails, error)
	UploadFiles(userID uint, filepaths []string) error
	FetchFilePathsByUserID(userID uint) ([]response.FetchUploadedDocumentPaths, error)
	RemoveFile(filePath string) error
	GetUserCount() (int, error)
}
