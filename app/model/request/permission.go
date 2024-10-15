package request

type RequestPermission struct {
	RoleID   uint   `json:"roleID"`
	Reason   string `json:"reason" binding:"required"`
	Date     string `json:"date" binding:"required"`
	FromTime string `json:"fromTime" binding:"required"`
	ToTime   string `json:"toTime" binding:"required"`
}

type UpdatePermissionStatus struct {
	IsApproved bool `json:"isApproved"`
}

type FetchUserPermissions struct {
	DepartmentMemberID uint `form:"departmentMemberID"`
	CommonRequestWithDateFilter
}
