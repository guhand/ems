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

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if user.DepartmentMemberID == nil {
		api_response.BadRequestError(c, "Department member not found")
		return
	}

	if err := h.leaveService.RequestLeave(*user.DepartmentMemberID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Leave requested successfully", nil)
}

func (h *LeaveHandler) FetchOwnLeaves(c *gin.Context) {
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

	if user.DepartmentMemberID == nil {
		api_response.BadRequestError(c, "Department member not found")
		return
	}

	data, err := h.leaveService.FetchOwnLeaves(*user.DepartmentMemberID, &filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Leaves fetched successfully", data)
}

func (h *LeaveHandler) FetchDepartmentMemberLeaves(c *gin.Context) {
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

	data, err := h.leaveService.FetchDepartmentMemberLeaves(*user.DepartmentID, &filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User leaves fetched successfully", data)
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

	api_response.Success(c, "Leave status updated successfully", nil)
}

func (h *LeaveHandler) FetchLeaves(c *gin.Context) {
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

	if user.DepartmentMemberID == nil {
		api_response.BadRequestError(c, "Department member not found")
		return
	}

	data, err := h.leaveService.FetchLeaves(&filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Leaves fetched successfully", data)
}

func (h *LeaveHandler) UpdateLeaveRequest(c *gin.Context) {
	var req request.RequestLeave

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Reason = utils.SqlParamValidator(req.Reason)

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
