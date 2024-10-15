package routes

import (
	"ems/api/middleware"
	"ems/infrastructure/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	departmentRepository := repository.NewDepartmentRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	leaveRepository := repository.NewLeaveRepository(db)
	permissionRepository := repository.NewPermissionRepository(db)
	noticeRepository := repository.NewNoticeRepository(db)

	middleware := middleware.NewMiddleware(userRepository)

	apiRoute := router.Group("api")

	RegisterAuthRoutes(apiRoute, userRepository, middleware)
	RegisterUserRoutes(apiRoute, userRepository, departmentRepository, leaveRepository, permissionRepository, middleware)
	RegisterDepartmentRoutes(apiRoute, departmentRepository, userRepository, middleware)
	RegisterRoleRoutes(apiRoute, roleRepository, middleware)
	RegisterLeaveRoute(apiRoute, leaveRepository, departmentRepository, userRepository, middleware)
	RegisterPermissionRoutes(apiRoute, permissionRepository, departmentRepository, userRepository, middleware)
	RegisterNoticeRoutes(apiRoute, noticeRepository, departmentRepository, middleware)
	RegisterDashboardRoutes(apiRoute, userRepository, departmentRepository, leaveRepository, permissionRepository, noticeRepository, middleware)
}
