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

type PermissionHandler struct {
	permissionService domain.PermissionService
}

func NewPermissionHandler(permissionService domain.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService}
}

func (h *PermissionHandler) RequestPermission(c *gin.Context) {
	var req request.RequestPermission

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Reason = utils.SqlParamValidator(req.Reason)
	req.Date = utils.SqlParamValidator(req.Date)
	req.FromTime = utils.SqlParamValidator(req.FromTime)
	req.ToTime = utils.SqlParamValidator(req.ToTime)

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if err := h.permissionService.RequestPermission(*user.DepartmentMemberID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Permission requested successfully", nil)
}

func (h *PermissionHandler) FetchOwnPermissions(c *gin.Context) {
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

	data, err := h.permissionService.FetchOwnPermissions(*user.DepartmentMemberID, &filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Permission requests fetched successfully", data)
}

func (h *PermissionHandler) FetchDepartmentMemberPermissions(c *gin.Context) {
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

	data, err := h.permissionService.FetchDepartmentMemberPermissions(*user.DepartmentID, &filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User permission requests fetched successfully", data)
}

func (h *PermissionHandler) UpdatePermissionStatus(c *gin.Context) {
	var req request.UpdatePermissionStatus

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.permissionService.UpdatePermissionStatus(uint(id), user.ID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Permission request status updated successfully", nil)
}

func (h *PermissionHandler) UpdatePermissionRequest(c *gin.Context) {
	var req request.RequestPermission

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Reason = utils.SqlParamValidator(req.Reason)
	req.Date = utils.SqlParamValidator(req.Date)
	req.FromTime = utils.SqlParamValidator(req.FromTime)
	req.ToTime = utils.SqlParamValidator(req.ToTime)

	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if err := h.permissionService.UpdatePermissionRequest(*user.DepartmentMemberID, uint(id), &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Permission request updated successfully", nil)
}

func (h *PermissionHandler) RemovePermissionRequest(c *gin.Context) {
	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.permissionService.RemovePermissionRequest(uint(id)); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Permission request removed successfully", nil)
}

func (h *PermissionHandler) FetchLeadAndHRPermissions(c *gin.Context) {
	var filters request.CommonRequestWithDateFilter

	if err := c.ShouldBindQuery(&filters); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	filters.Search = utils.SqlParamValidator(filters.Search)

	data, err := h.permissionService.FetchLeadAndHRPermissions(&filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "HR and Lead permission requests fetched successfully", data)
}

func (h *PermissionHandler) FetchUserPermissions(c *gin.Context) {
	var filters request.FetchUserPermissions

	if err := c.ShouldBindQuery(&filters); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	data, err := h.permissionService.FetchOwnPermissions(filters.DepartmentMemberID, &filters.CommonRequestWithDateFilter)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User permission requests fetched successfully", data)
}
