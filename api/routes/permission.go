package routes

import (
	"ems/api/middleware"
	"ems/app/handler"
	"ems/app/service"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

func RegisterPermissionRoutes(router *gin.RouterGroup, permissionRepository domain.PermissionRepository,
	departmentRepository domain.DepartmentRepository, userRepository domain.UserRepository,
	middleware *middleware.Middleware) {

	permissionService := service.NewPermissionService(permissionRepository, departmentRepository, userRepository)

	permissionHandler := handler.NewPermissionHandler(permissionService)

	userRoute := router.Group("permission", middleware.AuthMiddleware())
	{
		userRoute.POST("", permissionHandler.RequestPermission)
		userRoute.GET("", permissionHandler.FetchOwnPermissions)
		userRoute.PATCH(":id", permissionHandler.UpdatePermissionRequest)
		userRoute.DELETE(":id", permissionHandler.RemovePermissionRequest)
	}

	leadRoute := router.Group("lead/permission", middleware.DepartmentLeadMiddleware())
	{
		leadRoute.GET("", permissionHandler.FetchDepartmentMemberPermissions)
		leadRoute.PATCH(":id", permissionHandler.UpdatePermissionStatus)
	}

	managerRoute := router.Group("manager/permission", middleware.ManagerMiddleware())
	{
		managerRoute.GET("", permissionHandler.FetchLeadAndHRPermissions)
	}

	hrRoute := router.Group(`hr/permission`, middleware.HRAuthMiddleware())
	{
		hrRoute.GET("userPermission", permissionHandler.FetchUserPermissions)
	}
}
