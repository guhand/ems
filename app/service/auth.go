package service

import (
	apperror "ems/app/model/app_error"
	"ems/app/model/constant"
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/app/model/schema"
	"ems/domain"
	"ems/infrastructure/config"
	"ems/utils"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository domain.UserRepository
}

func NewAuthService(userRepository domain.UserRepository) domain.AuthService {
	return &authService{userRepository}
}

func (s *authService) Login(req *request.Login) (*response.FetchUserByEmail, error) {
	user, err := s.userRepository.GetUserByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.DataNotFoundError("email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	if (user.RoleID == uint(constant.Employee) || user.RoleID == uint(constant.DepartmentLead) ||
		user.RoleID == uint(constant.HR)) && user.DepartmentID == nil {
		return nil, fmt.Errorf("you are not assigned to any department, please contact HR")
	}

	token, err := utils.GenerateToken(int(user.ID))

	if err != nil {
		return nil, err
	}

	if err := s.userRepository.UpdateToken(user.ID, token); err != nil {
		return nil, err
	}

	user.Token = token

	return user, nil
}

func (s *authService) Logout(userID uint) error {
	isUserExistsByID, err := s.userRepository.GetUserByID(userID)

	if err != nil {
		return err
	}

	if isUserExistsByID == nil {
		return apperror.DataNotFoundError("user")
	}

	if err := s.userRepository.UpdateTokenStatus(userID); err != nil {
		return err
	}

	return nil
}

func (s *authService) SendForgotPasswordOtp(req *request.SendForgotPasswordOtp) error {
	isUserExistsByEmail, err := s.userRepository.GetUserByEmail(req.Email)

	if err != nil {
		return err
	}

	if isUserExistsByEmail == nil || isUserExistsByEmail.Email == "" {
		return apperror.DataNotFoundError("email")
	}

	otp := utils.GenerateOTP()

	userOtp := &schema.ForgotPasswordOtp{
		UserID: isUserExistsByEmail.ID,
		Email:  isUserExistsByEmail.Email,
		Otp:    otp}

	isOtpCreated, err := s.userRepository.CreateOTP(userOtp)

	if err != nil {
		return err
	}

	if isOtpCreated {
		if err := utils.SendFogotPasswordMail(req.Email, otp, time.Now()); err != nil {
			return err
		}
	}

	return nil
}

func (s *authService) VerifyForgotPasswordOtp(req *request.VerifyForgotPasswordOtp) (interface{}, error) {
	var response struct {
		Token string `json:"token"`
	}

	isUserExists, err := s.userRepository.GetUserByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	if isUserExists == nil || isUserExists.Email == "" {
		return nil, apperror.DataNotFoundError("email")
	}

	otpData, err := s.userRepository.GetOTPStatusByUserID(isUserExists.ID)

	if err != nil {
		return nil, err
	}

	if otpData.Otp != req.OTP || otpData.IsUsed {
		return nil, fmt.Errorf("incorrect otp")
	}

	expiresAt := otpData.CreatedAt.Add(time.Duration(config.Config.ForgotPasswordOTPValidity) * time.Minute)

	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("otp expired")
	}

	if err := s.userRepository.UpdateOTPStatus(isUserExists.ID); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(int(isUserExists.ID))

	if err != nil {
		return nil, err
	}

	if err := s.userRepository.UpdateToken(isUserExists.ID, token); err != nil {
		return nil, err
	}

	response.Token = token

	return response, nil
}
