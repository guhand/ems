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

	managerRoute := router.Group("hr/department", middleware.HRAuthMiddleware())

	{
		managerRoute.POST("", departmentHandler.CreateDepartment)
		managerRoute.GET("", departmentHandler.FetchDepartments)
		managerRoute.PATCH(":id", departmentHandler.UpdateDepartment)
		managerRoute.DELETE(":id", departmentHandler.RemoveDepartment)
		managerRoute.POST(":id/mapEmployees", departmentHandler.MappUsersToDepartment)
		managerRoute.GET(":id/users", departmentHandler.FetchDepartmentMembers)
		managerRoute.POST("unmapEmployee", departmentHandler.UnMapEmployee)
	}
}
