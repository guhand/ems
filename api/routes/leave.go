package routes

import (
	"ems/api/middleware"
	"ems/app/handler"
	"ems/app/service"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

func RegisterLeaveRoute(router *gin.RouterGroup, leaveRepository domain.LeaveRepository,
	departmentRepository domain.DepartmentRepository, userRepository domain.UserRepository,
	middleware *middleware.Middleware) {

	leaveService := service.NewLeaveService(leaveRepository, departmentRepository, userRepository)

	leaveHandler := handler.NewLeaveHandler(leaveService)

	userRoute := router.Group("leave", middleware.AuthMiddleware())
	{
		userRoute.POST("", leaveHandler.RequestLeave)
		userRoute.GET("", leaveHandler.FetchOwnLeaves)
		userRoute.PUT(":id", leaveHandler.UpdateLeaveRequest)
		userRoute.DELETE(":id", leaveHandler.RemoveLeaveRequest)
	}

	leadRoute := router.Group("lead/leave", middleware.DepartmentLeadMiddleware())
	{
		leadRoute.GET("", leaveHandler.FetchDepartmentMemberLeaves)
		leadRoute.PATCH(":id", leaveHandler.UpdateLeaveStatus)
	}

	managerRoute := router.Group("manager/leave", middleware.ManagerMiddleware())
	{
		managerRoute.GET("", leaveHandler.FetchLeadAndHRLeaves)
	}

	hrRoute := router.Group("hr/leave", middleware.HRAuthMiddleware())
	{
		hrRoute.GET("userLeave", leaveHandler.FetchUserLeaves)
	}
}
