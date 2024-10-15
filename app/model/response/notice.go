package response

import "time"

type FetchActiveUserNotices struct {
	DepartmentMemberID uint       `json:"departmentMemberID"`
	DepartmentMember   string     `json:"departmentMember" gorm:"column:departmentMember"`
	Role               string     `json:"role" gorm:"column:role"`
	CreatedAt          time.Time  `json:"createdAt"`
	Remarks            string     `json:"remarks"`
	IsApproved         bool       `json:"isApproved"`
	ApprovedBy         *string    `json:"approvedBy" gorm:"column:approvedBy"`
	NoticeEndDate      *time.Time `json:"noticeEndDate"`
	Count              int        `json:"-" gorm:"column:count"`
}
