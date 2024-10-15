package handler

import (
	"ems/api/api_response"
	"ems/api/middleware"
	"ems/app/model/request"
	"ems/domain"
	"ems/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	leaveService domain.LeaveService
}

func NewLeaveHandler(leaveService domain.LeaveService) *LeaveHandler {
	return &LeaveHandler{leaveService}
}

func (h *LeaveHandler) RequestLeave(c *gin.Context) {
	var req request.RequestLeave

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Reason = utils.SqlParamValidator(req.Reason)
	for i := range req.Dates {
		req.Dates[i].Date = utils.SqlParamValidator(req.Dates[i].Date)
	}

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if err := h.leaveService.RequestLeave(*user.DepartmentMemberID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Leave requested successfully", nil)
}

func (h *LeaveHandler) FetchOwnLeaves(c *gin.Context) {
	var filters request.CommonRequestWithDateFilter

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

	data, err := h.leaveService.FetchOwnLeaves(*user.DepartmentMemberID, &filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Leave requests fetched successfully", data)
}

func (h *LeaveHandler) FetchDepartmentMemberLeaves(c *gin.Context) {
	var filters request.CommonRequestWithDateFilter

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

	data, err := h.leaveService.FetchDepartmentMemberLeaves(*user.DepartmentID, &filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User leave requests fetched successfully", data)
}

func (h *LeaveHandler) UpdateLeaveStatus(c *gin.Context) {
	var req request.UpdateLeaveStatus

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.leaveService.UpdateLeaveStatus(uint(id), user.ID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Leave request status updated successfully", nil)
}

func (h *LeaveHandler) FetchLeadAndHRLeaves(c *gin.Context) {
	var filters request.CommonRequestWithDateFilter

	if err := c.ShouldBindQuery(&filters); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	filters.Search = utils.SqlParamValidator(filters.Search)

	data, err := h.leaveService.FetchLeadAndHRLeaves(&filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "HR and Lead leave requests fetched successfully", data)
}

func (h *LeaveHandler) UpdateLeaveRequest(c *gin.Context) {
	var req request.RequestLeave

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Reason = utils.SqlParamValidator(req.Reason)

	for i := range req.Dates {
		req.Dates[i].Date = utils.SqlParamValidator(req.Dates[i].Date)
	}

	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.leaveService.UpdateLeaveRequest(uint(id), &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Leave request updated successfully", nil)
}

func (h *LeaveHandler) RemoveLeaveRequest(c *gin.Context) {
	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.leaveService.RemoveLeaveRequest(uint(id)); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Leave request removed successfully", nil)
}

func (h *LeaveHandler) FetchUserLeaves(c *gin.Context) {
	var filters request.FetchUserLeaves

	if err := c.ShouldBindQuery(&filters); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	filters.Search = utils.SqlParamValidator(filters.Search)

	data, err := h.leaveService.FetchOwnLeaves(filters.DepartmentMemberID, &filters.CommonRequestWithDateFilter)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User leave requests fetched successfully", data)
}
