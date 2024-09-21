package domain

import (
	"ems/app/model/request"
	"ems/app/model/response"
)

type NoticeService interface {
	CreateNotice(departmentMemberID uint, req *request.CreateNotice) error
	FetchActiveUserNotices(filters *request.CommonRequest) ([]response.FetchActiveUserNotices, error)
	FetchNotice(departmentMemberID uint) (*response.FetchActiveUserNotices, error)
	ApproveNotice(departmentMemberID, approvedBy uint, req *request.ApproveNotice) error
}

type NoticeRepository interface {
	CreateNotice(departmentMemberID uint, req *request.CreateNotice) error
	FetchActiveUserNotices(filters *request.CommonRequest) ([]response.FetchActiveUserNotices, error)
	GetNoticeUserCount() (int, error)
	GetNoticeUserCountByDepartment(departmentID uint) (int, error)
	FetchNotice(departmentMemberID uint) (*response.FetchActiveUserNotices, error)
	ApproveNotice(departmentMemberID, approvedBy uint, req *request.ApproveNotice) error
}
