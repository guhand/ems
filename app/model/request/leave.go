package request

type RequestLeave struct {
	RoleID uint   `json:"roleID"`
	Reason string `json:"reason" binding:"required"`
	Dates  []Date `json:"dates"`
}

type Date struct {
	Date        string `json:"date" binding:"required"`
	IsFullDay   bool   `json:"isFullDay"`
	SessionType uint   `json:"sessionType"`
}

type UpdateLeaveStatus struct {
	IsApproved bool `json:"isApproved"`
}

type FetchUserLeaves struct {
	DepartmentMemberID uint `form:"departmentMemberID"`
	CommonRequestWithDateFilter
}
