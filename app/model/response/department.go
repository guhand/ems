package response

import "time"

type FetchDepartments struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	LeadId          *uint     `json:"leadID" gorm:"column:leadID"`
	LeadName        *string   `json:"leadName" gorm:"column:leadName"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	IsActive        bool      `json:"isActive"`
	DepartmentCount uint      `json:"departmentCount" gorm:"column:departmentCount"`
}

type FetchDepartmentMembers struct {
	UserID                uint      `json:"userID" gorm:"column:userID"`
	UserName              string    `json:"userName" gorm:"column:userName"`
	Role                  string    `json:"role" gorm:"column:role"`
	CreatedAt             time.Time `json:"createdAt"`
	DepartmentMemberCount uint      `json:"departmentMemberCount" gorm:"column:departmentMemberCount"`
}

type FetchDepartmenLead struct {
	LeadId   uint   `json:"leadID" gorm:"column:leadID"`
	LeadName string `json:"leadName" gorm:"column:leadName"`
}
