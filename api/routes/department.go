package routes

import (
	"ems/api/middleware"
	"ems/app/handler"
	"ems/app/service"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

func RegisterDepartmentRoutes(router *gin.RouterGroup,
	departmentRepository domain.DepartmentRepository, userRepository domain.UserRepository,
	middleware *middleware.Middleware) {

	departmentService := service.NewDepartmentService(departmentRepository, userRepository)

	departmentHandler := handler.NewDepartmentHandler(departmentService)

	hrRoute := router.Group("department", middleware.HRAuthMiddleware())

	{
		hrRoute.POST("", departmentHandler.CreateDepartment)
		hrRoute.GET("", departmentHandler.FetchDepartments)
		hrRoute.PATCH(":id", departmentHandler.UpdateDepartment)
		hrRoute.DELETE(":id", departmentHandler.RemoveDepartment)
		hrRoute.POST(":id/mapUsers", departmentHandler.MappUsersToDepartment)
		hrRoute.GET(":id/users", departmentHandler.FetchDepartmentMembers)
		hrRoute.POST("unmapUser", departmentHandler.UnMapUser)
	}
}
