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

	employeeRoute := router.Group("employee/permission", middleware.AuthMiddleware())
	{
		employeeRoute.POST("", permissionHandler.RequestPermission)
		employeeRoute.GET("", permissionHandler.FetchPermissions)
		employeeRoute.PATCH(":id", permissionHandler.UpdatePermissionRequest)
		employeeRoute.DELETE(":id", permissionHandler.RemovePermissionRequest)
	}

	leadRoute := router.Group("lead/permission", middleware.DepartmentLeadMiddleware())
	{
		leadRoute.GET("", permissionHandler.FetchUserPermissions)
		leadRoute.PATCH(":id", permissionHandler.UpdatePermissionStatus)
	}
}
