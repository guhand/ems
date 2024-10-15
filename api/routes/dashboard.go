package routes

import (
	"ems/api/middleware"
	"ems/app/handler"
	"ems/app/service"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

func RegisterDashboardRoutes(router *gin.RouterGroup, userRepository domain.UserRepository,
	departRepository domain.DepartmentRepository, leaveRepository domain.LeaveRepository,
	permissionRepository domain.PermissionRepository, noticeRepository domain.NoticeRepository,
	middleware *middleware.Middleware) {

	dashboardService := service.NewDashboardService(userRepository, departRepository, leaveRepository,
		permissionRepository, noticeRepository)

	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	dashboardRoute := router.Group("dashboard", middleware.AuthMiddleware())
	{
		dashboardRoute.GET("", dashboardHandler.FetchDashboardCounts)
	}
}
