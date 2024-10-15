package response

import "time"

type FetchUserByEmail struct {
	ID                 uint      `json:"id"`
	FirstName          string    `json:"firstName"`
	LastName           string    `json:"lastName"`
	Code               string    `json:"code"`
	Email              string    `json:"email"`
	Mobile             string    `json:"mobile"`
	Token              string    `json:"token"`
	Password           string    `json:"-"`
	ManagerID          *uint     `json:"managerID,omitempty" gorm:"column:managerID"`
	Manager            *string   `json:"manager,omitempty" gorm:"column:manager"`
	RoleID             uint      `json:"roleID" gorm:"column:roleID"`
	Role               string    `json:"role" gorm:"column:roleName"`
	DepartmentID       *uint     `json:"departmentID,omitempty" gorm:"column:departmentID"`
	Department         *string   `json:"department,omitempty" gorm:"column:department"`
	DepartmentMemberID *uint     `json:"departmentMemberID,omitempty" gorm:"column:departmentMemberID"`
	LeadID             *uint     `json:"leadID,omitempty" gorm:"column:leadID"`
	Lead               *string   `json:"lead,omitempty" gorm:"column:lead"`
	CreatedAt          time.Time `json:"createdAt"`
	IsActive           bool      `json:"isActive"`
}

type FetchUserByID struct {
	ID                 uint    `json:"userID" gorm:"column:userID"`
	Token              *string `json:"token" gorm:"column:token"`
	RoleID             uint    `json:"roleID" gorm:"column:roleID"`
	DepartmentID       *uint   `json:"departmentID" gorm:"column:departmentID"`
	DepartmentMemberID *uint   `json:"departmentMemberID" gorm:"column:departmentMemberID"`
}

type FetchUsers struct {
	ID                 uint      `json:"id"`
	FirstName          string    `json:"firstName" gorm:"column:firstName"`
	LastName           string    `json:"lastName" gorm:"column:lastName"`
	Email              string    `json:"email"`
	Mobile             string    `json:"mobile"`
	Code               string    `json:"code"`
	DepartmentMemberID *uint     `json:"departmentMemberID,omitempty" gorm:"column:departmentMemberID"`
	RoleID             uint      `json:"roleID" gorm:"column:roleID"`
	RoleName           string    `json:"roleName" gorm:"column:roleName"`
	DepartmentID       *uint     `json:"departmentID,omitempty" gorm:"column:departmentID"`
	Department         string    `json:"department" gorm:"column:department"`
	DateOfJoining      *string   `json:"dateOfJoining"`
	Experience         *string   `json:"experience"`
	Designation        *string   `json:"designation"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	IsActive           bool      `json:"isActive"`
	UserCount          uint      `json:"userCount" gorm:"column:userCount"`
}

type FetchLastUserCode struct {
	Code string `json:"code"`
}

type FetchUnmappedUsers struct {
	ID   uint   `json:"id"`
	Name string `json:"name" gorm:"column:name"`
}

type FetchUserDetails struct {
	UserID            uint    `json:"userID" gorm:"column:userID"`
	DateOfJoining     *string `json:"dateOfJoining"`
	Experience        *string `json:"experience"`
	Designation       *string `json:"designation"`
	DOB               *string `json:"dob"`
	PanNumber         *string `json:"panNumber"`
	AadharNumber      *string `json:"aadharNumber"`
	BankAccountNumber *string `json:"bankAccountNumber"`
	IfscCode          *string `json:"ifscCode"`
	City              *string `json:"city"`
	Address           *string `json:"address"`
	Degree            *string `json:"degree"`
	College           *string `json:"college"`
}

type FetchUploadedDocumentPaths struct {
	UserID    uint   `json:"userID"`
	FilePaths string `json:"filePaths" gorm:"column:filePaths"`
}

type FetchDepartmentUserCountAndRoleID struct {
	Count        int `gorm:"column:count"`
	RoleID       int `gorm:"column:roleID"`
	DepartmentID int `gorm:"column:departmentID"`
}
