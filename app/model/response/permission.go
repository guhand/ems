package response

import "time"

type FetchUserPermissions struct {
	ID                  uint       `json:"id"`
	DepartmentMemberID  uint       `json:"departmentMemberID" gorm:"column:departmentMemberID"`
	DepartmentMember    string     `json:"departmentMember" gorm:"column:departmentMember"`
	Role                *string    `json:"role,omitempty" gorm:"column:role"`
	Reason              string     `json:"reason"`
	Date                string     `json:"date" gorm:"column:date"`
	FromTime            string     `json:"fromTime" gorm:"column:fromTime"`
	ToTime              string     `json:"toTime" gorm:"column:toTime"`
	IsApproved          *bool      `json:"isApproved" gorm:"column:isApproved"`
	ApprovedAt          *time.Time `json:"approvedAt" gorm:"column:approvedAt"`
	ApprovedBy          *string    `json:"approvedBy" gorm:"column:approvedBy"`
	IsActive            bool       `json:"isActive"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	Count               int        `json:"-" gorm:"column:count"`
	UserPermissionCount int        `json:"userPermissionCount" gorm:"-"`
}
