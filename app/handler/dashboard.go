package handler

import (
	"ems/api/api_response"
	"ems/api/middleware"
	"ems/app/model/request"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService domain.DashboardService
}

func NewDashboardHandler(dashboardService domain.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService}
}

func (h *DashboardHandler) FetchDashboardCounts(c *gin.Context) {
	var req request.DateFilters

	if err := c.ShouldBindQuery(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	user, err := middleware.GetUserClaims(c)
	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	var (
		departmentID       uint
		departmentMemberID uint
	)

	if user.DepartmentID != nil {
		departmentID = *user.DepartmentID
	}

	if user.DepartmentMemberID != nil {
		departmentMemberID = *user.DepartmentMemberID
	}

	data, err := h.dashboardService.FetchDashboardCounts(user.RoleID, departmentMemberID, departmentID, &req)
	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Dashboard counts fetched successfully", data)
}
