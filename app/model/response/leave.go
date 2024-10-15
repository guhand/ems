package response

import "time"

type FetchLeaves struct {
	ID                 uint       `json:"id"`
	DepartmentMemberID uint       `json:"departmentMemberID" gorm:"column:departmentMemberID"`
	DepartmentMember   string     `json:"departmentMember" gorm:"column:departmentMember"`
	Role               *string    `json:"role,omitempty" gorm:"column:role"`
	Reason             string     `json:"reason"`
	Dates              string     `json:"dates" gorm:"column:dates"`
	IsFullDays         string     `json:"isFullDays" gorm:"column:isFullDays"`
	SessionTypes       string     `json:"sessionTypes" gorm:"column:sessionTypes"`
	IsApproved         *bool      `json:"isApproved" gorm:"column:isApproved"`
	ApprovedAt         *time.Time `json:"approvedAt" gorm:"column:approvedAt"`
	ApprovedBy         *string    `json:"approvedBy" gorm:"column:approvedBy"`
	IsActive           bool       `json:"isActive"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	Count              uint       `json:"-" gorm:"column:count"`
	UserLeaveCount     float64    `json:"userLeaveCount" gorm:"-"`
}
