package request

import "time"

type CreateUser struct {
	RoleID    uint   `json:"roleID" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Mobile    string `json:"mobile" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type UpdateUser struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Mobile    string `json:"mobile" binding:"required"`
}

type FetchUnmappedLeadUserIncludeUserID struct {
	DepartmentID uint `json:"departmentID" binding:"required"`
	UserID       uint `json:"userID" binding:"required"`
}

type UpdateUserDetails struct {
	UserID            uint      `json:"userID" binding:"required"`
	DateOfJoining     time.Time `json:"dateOfJoining" binding:"required"`
	Experience        string    `json:"experience" binding:"required"`
	Designation       string    `json:"designation" binding:"required"`
	DOB               time.Time `json:"dob" binding:"required"`
	PanNumber         string    `json:"panNumber" binding:"required"`
	AadharNumber      string    `json:"aadharNumber" binding:"required"`
	BankAccountNumber string    `json:"bankAccountNumber" binding:"required"`
	IfscCode          string    `json:"ifscCode" binding:"required"`
	City              string    `json:"city" binding:"required"`
	Address           string    `json:"address" binding:"required"`
	Degree            string    `json:"degree" binding:"required"`
	College           string    `json:"college" binding:"required"`
}

type FetchUserDetails struct {
	UserID uint `form:"userID" binding:"required"`
}

type FetchUnmappedUsersByDepartmentID struct {
	DepartmentID uint `json:"departmentID" binding:"required"`
}

type FetchUsers struct {
	CommonRequest
	DepartmentID uint `form:"departmentID"`
	RoleID       uint `form:"roleID"`
}
