package request

type ApplyNotice struct {
	Remarks string `json:"remarks" binding:"required"`
}

type ApproveNotice struct {
	DepartmentMemberID int `json:"departmentMemberID" binding:"required"`
	ServeDays          int `json:"serveDays" binding:"required"`
}
