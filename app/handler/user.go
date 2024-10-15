package handler

import (
	"ems/api/api_response"
	"ems/api/middleware"
	"ems/app/model/request"
	"ems/domain"
	"ems/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.CreateUser

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.FirstName = utils.SqlParamValidator(req.FirstName)
	req.LastName = utils.SqlParamValidator(req.LastName)
	req.Code = utils.SqlParamValidator(req.Code)
	req.Email = utils.SqlParamValidator(req.Email)
	req.Mobile = utils.SqlParamValidator(req.Mobile)
	req.Password = utils.SqlParamValidator(req.Password)

	if err := h.userService.CreateUser(&req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User created successfully", nil)
}

func (h *UserHandler) FetchUsers(c *gin.Context) {
	var filters request.FetchUsers

	if err := c.ShouldBindQuery(&filters); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	filters.Search = utils.SqlParamValidator(filters.Search)

	data, err := h.userService.FetchUsers(&filters)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Users fetched successfully", data)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req request.UpdateUser

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.FirstName = utils.SqlParamValidator(req.FirstName)
	req.LastName = utils.SqlParamValidator(req.LastName)
	req.Code = utils.SqlParamValidator(req.Code)
	req.Email = utils.SqlParamValidator(req.Email)
	req.Mobile = utils.SqlParamValidator(req.Mobile)

	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.userService.UpdateUser(uint(id), &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User updated successfully", nil)
}

func (h *UserHandler) FetchLastUserCode(c *gin.Context) {
	data, err := h.userService.FetchLastUserCode()

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Last user code fetched successfully", data)
}

func (h *UserHandler) FetchUnmappedLeadUsers(c *gin.Context) {

	data, err := h.userService.FetchUnmappedLeadUsers()

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Unmapped department lead users fetched successfully", data)
}

func (h *UserHandler) FetchUnmappedLeadUserIncludeUserID(c *gin.Context) {
	var req request.FetchUnmappedLeadUserIncludeUserID

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	data, err := h.userService.FetchUnmappedLeadUserIncludeUserID(&req)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Unmapped department lead users include id fetched successfully", data)
}

func (h *UserHandler) FetchUnmappedUsers(c *gin.Context) {
	var req request.FetchUnmappedUsersByDepartmentID

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	data, err := h.userService.FetchUnmappedUsers(&req)

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Unmapped users fetched successfully", data)
}

func (h *UserHandler) RemoveUser(c *gin.Context) {
	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	if err := h.userService.RemoveUser(uint(id)); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User removed successfully", nil)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req request.ResetPassword

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Password = utils.SqlParamValidator(req.Password)

	user, err := middleware.GetUserClaims(c)

	if err != nil {
		api_response.UnauthorizedError(c, err.Error())
		return
	}

	if err := h.userService.ResetPassword(user.ID, &req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Password reset successfully", nil)
}

func (h *UserHandler) UpdateUserDetails(c *gin.Context) {
	var req request.UpdateUserDetails

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.PanNumber = utils.SqlParamValidator(req.PanNumber)
	req.AadharNumber = utils.SqlParamValidator(req.AadharNumber)
	req.BankAccountNumber = utils.SqlParamValidator(req.BankAccountNumber)
	req.IfscCode = utils.SqlParamValidator(req.IfscCode)
	req.City = utils.SqlParamValidator(req.City)
	req.Address = utils.SqlParamValidator(req.Address)
	req.Degree = utils.SqlParamValidator(req.Degree)
	req.College = utils.SqlParamValidator(req.College)

	if err := h.userService.UpdateUserDetails(&req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User details updated successfully", nil)
}

func (h *UserHandler) FetchUserDetails(c *gin.Context) {
	var req request.FetchUserDetails

	if err := c.ShouldBindQuery(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	data, err := h.userService.FetchUserDetails(&req)
	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User details fetched successfully", data)
}

func (h *UserHandler) UploadFiles(c *gin.Context) {
	userIdForm := c.PostForm("userID")
	if userIdForm == "" {
		api_response.BadRequestError(c, "userId is required")
		return
	}

	userID, err := strconv.Atoi(userIdForm)
	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		api_response.BadRequestError(c, "No files provided")
		return
	}

	baseDir := "./uploads/" + "user-" + userIdForm

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}
	}

	var filepaths []string

	for _, file := range files {
		if filepath.Ext(file.Filename) != ".pdf" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
			return
		}

		filePath := filepath.Join(baseDir, file.Filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload %s", file.Filename)})
			return
		}

		filepaths = append(filepaths, filePath)
	}

	if err := h.userService.UploadFiles(uint(userID), filepaths); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "User documents uploaded successfully", nil)
}

func (h *UserHandler) FetchFilePathsByUserID(c *gin.Context) {
	userIdForm := c.Query("userID")
	if userIdForm == "" {
		api_response.BadRequestError(c, "userID is required")
		return
	}

	userID, err := strconv.Atoi(userIdForm)
	if err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	data, err := h.userService.FetchFilePathsByUserID(uint(userID))
	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "File paths fetched successfully", data)
}

func (h *UserHandler) FetchFile(c *gin.Context) {
	filePath := c.Query("path")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		api_response.BadRequestError(c, "File not found")
		return
	}

	c.File(filePath)
}

func (h *UserHandler) DeleteFile(c *gin.Context) {
	filePath := c.Query("path")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		api_response.BadRequestError(c, "File not found")
		return
	}

	err := os.Remove(filePath)
	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	if err := h.userService.RemoveFile(filePath); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "File deleted successfully", nil)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req request.ChangePassword

	if err := c.ShouldBindJSON(&req); err != nil {
		api_response.BadRequestError(c, err.Error())
		return
	}

	req.Email = utils.SqlParamValidator(req.Email)
	req.OldPassword = utils.SqlParamValidator(req.OldPassword)
	req.NewPassword = utils.SqlParamValidator(req.NewPassword)

	if err := h.userService.ChangePassword(&req); err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Password changed successfully", nil)
}

func (h *UserHandler) FetchUnmappedHRUsers(c *gin.Context) {

	data, err := h.userService.FetchUnmappedHRUsers()

	if err != nil {
		api_response.InternalServerError(c, err.Error())
		return
	}

	api_response.Success(c, "Unmapped hr users fetched successfully", data)
}
