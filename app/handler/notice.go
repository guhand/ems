package handler

import (
	"ems/api/api_response"
	"ems/api/middleware"
	"ems/app/model/request"
	"ems/domain"
	"ems/utils"

	"github.com/gin-gonic/gin"
)

type NoticeHandler struct {
	noticeService domain.NoticeService
}

func NewNoticeHandler(noticeService domain.NoticeService) *NoticeHandler {
	return &NoticeHandler{noticeService}
}

func (h *NoticeHandler) ApplyNotice(c *gin.Context) {
	var req request.ApplyNotice

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if user.DepartmentMemberID == nil {
		api_response.BadRequestError(c, "you are not assigned to any department. Kindly contact the manager")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Remarks = utils.SqlParamValidator(req.Remarks)

	if err := h.noticeService.ApplyNotice(*user.DepartmentMemberID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Notice applied successfully", nil)
}

func (h *NoticeHandler) FetchActiveUserNotices(c *gin.Context) {
	var filters request.CommonRequest

	if err := c.ShouldBindQuery(&filters); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	filters.Search = utils.SqlParamValidator(filters.Search)

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	data, err := h.noticeService.FetchActiveUserNotices(user.RoleID, &filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return

	}

	api_response.Success(c, "User notices fetched successfully", data)
}

func (h *NoticeHandler) FetchNotice(c *gin.Context) {
	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if user.DepartmentMemberID == nil {
		api_response.BadRequestError(c, "you are not assigned to any department. Kindly contact the manager")
		return
	}

	data, err := h.noticeService.FetchNotice(*user.DepartmentMemberID)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User notice fetched successfully", data)
}

func (h *NoticeHandler) ApproveNotice(c *gin.Context) {
	var req request.ApproveNotice

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if err := h.noticeService.ApproveNotice(user.ID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User notice approved successfully", nil)
}
