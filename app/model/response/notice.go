package response

import "time"

type FetchActiveUserNotices struct {
	DepartmentMemberID uint       `json:"departmentMemberID"`
	DepartmentMember   string     `json:"departmentMember" gorm:"column:departmentMember"`
	CreatedAt          time.Time  `json:"createdAt"`
	Remarks            string     `json:"remarks"`
	IsApproved         bool       `json:"isApproved"`
	ApprovedBy         *uint      `json:"approvedBy" gorm:"column:approvedBy"`
	NoticeEndDate      *time.Time `json:"noticeEndDate"`
}
