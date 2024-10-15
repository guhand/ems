package handler

import (
	"ems/api/api_response"
	"ems/domain"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService domain.RoleService
}

func NewRoleHandler(roleService domain.RoleService) *RoleHandler {
	return &RoleHandler{roleService}
}

func (h *RoleHandler) FetchRoles(c *gin.Context) {

	data, err := h.roleService.FetchRoles()

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Roles fetched successfully", data)
}
