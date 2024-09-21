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

	employeeRoute := router.Group("employee/leave", middleware.AuthMiddleware())
	{
		employeeRoute.POST("", leaveHandler.RequestLeave)
		employeeRoute.GET("", leaveHandler.FetchOwnLeaves)
		employeeRoute.PUT(":id", leaveHandler.UpdateLeaveRequest)
		employeeRoute.DELETE(":id", leaveHandler.RemoveLeaveRequest)
	}

	leadRoute := router.Group("lead/leave", middleware.DepartmentLeadMiddleware())
	{
		leadRoute.GET("", leaveHandler.FetchDepartmentMemberLeaves)
		leadRoute.PATCH(":id", leaveHandler.UpdateLeaveStatus)
	}

	managerRoute := router.Group("manager/leave", middleware.ManagerMiddleware())
	{
		managerRoute.GET("", leaveHandler.FetchLeaves)
	}
}
