package routes

import (
	"ems/api/middleware"
	"ems/app/handler"
	"ems/app/service"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, userRepository domain.UserRepository,
	middleware *middleware.Middleware) {

	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

	authRoute := router.Group("auth")
	{
		authRoute.POST("login", authHandler.Login)
		authRoute.POST("logout", middleware.AuthMiddleware(), authHandler.Logout)
	}

	forgotPasswordRoute := router.Group("forgotPassword")
	{
		forgotPasswordRoute.POST("sendOtp", authHandler.SendForgotPasswordOtp)
		forgotPasswordRoute.POST("verifyOtp", authHandler.VerifyForgotPasswordOtp)
	}
}
