package domain

import (
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/utils"
)

type NoticeService interface {
	ApplyNotice(departmentMemberID uint, req *request.ApplyNotice) error
	FetchActiveUserNotices(roleID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	FetchNotice(departmentMemberID uint) (*response.FetchActiveUserNotices, error)
	ApproveNotice(approvedBy uint, req *request.ApproveNotice) error
}

type NoticeRepository interface {
	ApplyNotice(departmentMemberID uint, req *request.ApplyNotice) error
	FetchActiveUserNotices(roleID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error)
	GetNoticeUserCount() (int, error)
	GetNoticeUserCountByDepartment(departmentID uint) (int, error)
	FetchNotice(departmentMemberID uint) (*response.FetchActiveUserNotices, error)
	ApproveNotice(departmentMemberID, approvedBy uint, req *request.ApproveNotice) error
	IsApproveExistsByUser(departmentMemberID uint) (bool, error)
}
