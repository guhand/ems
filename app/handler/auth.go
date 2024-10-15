package handler

import (
	"ems/api/api_response"
	"ems/api/middleware"
	"ems/app/model/request"
	"ems/domain"
	"ems/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.Login

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Email = utils.SqlParamValidator(req.Email)
	req.Password = utils.SqlParamValidator(req.Password)

	data, err := h.authService.Login(&req)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Login successful", data)
}

func (h *AuthHandler) Logout(c *gin.Context) {

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if err := h.authService.Logout(user.ID); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Logout successful", nil)
}

func (h *AuthHandler) SendForgotPasswordOtp(c *gin.Context) {
	var req request.SendForgotPasswordOtp

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.authService.SendForgotPasswordOtp(&req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "OTP sent successfully", nil)
}

func (h *AuthHandler) VerifyForgotPasswordOtp(c *gin.Context) {
	var req request.VerifyForgotPasswordOtp

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Email = utils.SqlParamValidator(req.Email)
	req.OTP = utils.SqlParamValidator(req.OTP)

	data, err := h.authService.VerifyForgotPasswordOtp(&req)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "OTP verified successfully", data)
}
