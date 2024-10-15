package handler

import (
	"ems/api/api_response"
	"ems/app/model/request"
	"ems/domain"
	"ems/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	departmentService domain.DepartmentService
}

func NewDepartmentHandler(departmentService domain.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{departmentService}
}

func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req request.CreateDepartment

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Name = utils.SqlParamValidator(req.Name)

	if err := h.departmentService.CreateDepartment(&req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Department created successfully", nil)
}

func (h *DepartmentHandler) FetchDepartments(c *gin.Context) {
	var filters request.CommonRequest

	if err := c.ShouldBindQuery(&filters); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	filters.Search = utils.SqlParamValidator(filters.Search)

	data, err := h.departmentService.FetchDepartments(&filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Departments fetched successfully", data)
}

func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	var req request.UpdateDepartment

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Name = utils.SqlParamValidator(req.Name)

	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.departmentService.UpdateDepartment(uint(id), &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Department updated successfully", nil)
}

func (h *DepartmentHandler) RemoveDepartment(c *gin.Context) {
	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.departmentService.RemoveDepartment(uint(id)); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Department removed successfully", nil)
}

func (h *DepartmentHandler) MappUsersToDepartment(c *gin.Context) {
	var req request.MappUsersToDepartment

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

	if err := h.departmentService.MappUsersToDepartment(uint(id), &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Users mapped to the department successfully", nil)
}

func (h *DepartmentHandler) FetchDepartmentMembers(c *gin.Context) {
	var filters request.CommonRequest

	if err := c.ShouldBindQuery(&filters); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	filters.Search = utils.SqlParamValidator(filters.Search)

	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	data, err := h.departmentService.FetchDepartmentMembers(uint(id), &filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Department members fetched successfully", data)
}

func (h *DepartmentHandler) UnMapUser(c *gin.Context) {
	var req request.UnMapUser

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.departmentService.UnMapUser(&req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User unmapped from the department successfully", nil)
}
