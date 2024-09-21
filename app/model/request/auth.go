package request

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SendForgotPasswordOtp struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyForgotPasswordOtp struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,numeric,len=6"`
}

type ResetPassword struct {
	Password string `json:"password" binding:"required"`
}
