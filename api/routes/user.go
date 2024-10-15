package routes

import (
	"ems/api/middleware"
	"ems/app/handler"
	"ems/app/service"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, userRepository domain.UserRepository,
	departmentRepository domain.DepartmentRepository, leaveRepository domain.LeaveRepository,
	permissionRepository domain.PermissionRepository,
	middleware *middleware.Middleware) {

	userService := service.NewUserService(userRepository, departmentRepository, leaveRepository, permissionRepository)
	userHandler := handler.NewUserHandler(userService)

	hrRoute := router.Group("hr/user", middleware.HRAuthMiddleware())
	{
		hrRoute.POST("", userHandler.CreateUser)
		hrRoute.GET("", userHandler.FetchUsers)
		hrRoute.PATCH(":id", userHandler.UpdateUser)
		hrRoute.DELETE(":id", userHandler.RemoveUser)
		hrRoute.POST("details", userHandler.UpdateUserDetails)
		hrRoute.GET("lastUserCode", userHandler.FetchLastUserCode)
		hrRoute.GET("unmappedLeads", userHandler.FetchUnmappedLeadUsers)
		hrRoute.GET("unmappedHrs", userHandler.FetchUnmappedHRUsers)
		hrRoute.POST("unmappedLeadsIncludeID", userHandler.FetchUnmappedLeadUserIncludeUserID)
		hrRoute.POST("fetchUnmappedUsers", userHandler.FetchUnmappedUsers)
		hrRoute.POST("uploadFiles", userHandler.UploadFiles)
		hrRoute.GET("/files", userHandler.FetchFilePathsByUserID)
		hrRoute.GET("/file", userHandler.FetchFile)
		hrRoute.DELETE("/file", userHandler.DeleteFile)
	}

	userRoute := router.Group("user", middleware.AuthMiddleware())
	{
		userRoute.GET("details", userHandler.FetchUserDetails)
		userRoute.POST("resetPassword", userHandler.ResetPassword)
		userRoute.POST("changePassword", userHandler.ChangePassword)
	}
}
