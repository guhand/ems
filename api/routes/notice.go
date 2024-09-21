package routes

import (
	"ems/api/middleware"
	"ems/app/handler"
	"ems/app/service"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

func RegisterNoticeRoutes(router *gin.RouterGroup, noticeRepository domain.NoticeRepository,
	departmentRepository domain.DepartmentRepository,
	middleware *middleware.Middleware) {

	noticeService := service.NewNoticeService(noticeRepository, departmentRepository)

	noticeHandler := handler.NewNoticeHandler(noticeService)

	employeeRoute := router.Group("employee/notice", middleware.AuthMiddleware())
	{
		employeeRoute.POST("", noticeHandler.CreateNotice)
		employeeRoute.GET("", noticeHandler.FetchNotice)
	}

	hrRoute := router.Group("hr/notice", middleware.HRAuthMiddleware())
	{
		hrRoute.GET("", noticeHandler.FetchActiveUserNotices)
		hrRoute.POST("", noticeHandler.ApproveNotice)
	}
}
