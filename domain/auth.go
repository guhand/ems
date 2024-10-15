package domain

import (
	"ems/app/model/request"
	"ems/app/model/response"
)

type AuthService interface {
	Login(req *request.Login) (*response.FetchUserByEmail, error)
	Logout(userID uint) error
	SendForgotPasswordOtp(req *request.SendForgotPasswordOtp) error
	VerifyForgotPasswordOtp(req *request.VerifyForgotPasswordOtp) (interface{}, error)
}
