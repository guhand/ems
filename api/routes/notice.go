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

	userRoute := router.Group("notice", middleware.AuthMiddleware())
	{
		userRoute.POST("", noticeHandler.ApplyNotice)
		userRoute.GET("", noticeHandler.FetchNotice)
	}

	hrRoute := router.Group("hr/notice", middleware.HRAuthMiddleware())
	{
		hrRoute.GET("", noticeHandler.FetchActiveUserNotices)
		hrRoute.POST("", noticeHandler.ApproveNotice)
	}
}
