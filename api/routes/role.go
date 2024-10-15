package routes

import (
	"ems/api/middleware"
	"ems/app/handler"
	"ems/app/service"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(router *gin.RouterGroup, roleRepository domain.RoleRepository,
	middleware *middleware.Middleware) {

	roleService := service.NewRoleService(roleRepository)

	roleHandler := handler.NewRoleHandler(roleService)

	roleRoute := router.Group("role", middleware.HRAuthMiddleware())
	{
		roleRoute.GET("", roleHandler.FetchRoles)
	}
}
