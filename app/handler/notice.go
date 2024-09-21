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

func (h *NoticeHandler) CreateNotice(c *gin.Context) {
	var req request.CreateNotice

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Remarks = utils.SqlParamValidator(req.Remarks)

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if err := h.noticeService.CreateNotice(*user.DepartmentMemberID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Notice created successfully", nil)
}

func (h *NoticeHandler) FetchActiveUserNotices(c *gin.Context) {
	var filters request.CommonRequest

	if err := c.ShouldBindQuery(&filters); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	filters.Search = utils.SqlParamValidator(filters.Search)

	data, err := h.noticeService.FetchActiveUserNotices(&filters)

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

	if err := h.noticeService.ApproveNotice(*user.DepartmentMemberID, user.ID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Notice approved successfully", nil)
}
