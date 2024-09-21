package repository

import (
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/domain"
	"time"

	"gorm.io/gorm"
)

type noticeRepository struct {
	db *gorm.DB
}

func NewNoticeRepository(db *gorm.DB) domain.NoticeRepository {
	return &noticeRepository{db}
}

func (r *noticeRepository) CreateNotice(departmentMemberID uint, req *request.CreateNotice) error {
	return r.db.Exec(`
			INSERT INTO UserNotice
			(CreatedAt, UpdatedAt, DepartmentMemberID, Remarks)
			VALUES(?, ?, ?, ?)`,
		time.Now(), time.Now(), departmentMemberID, req.Remarks).Error
}

func (r *noticeRepository) FetchActiveUserNotices(filters *request.CommonRequest) ([]response.FetchActiveUserNotices, error) {
	var data []response.FetchActiveUserNotices

	now := time.Now()
	itemsPerPage := 10

	if err := r.db.Raw(`
		SELECT un.DepartmentMemberID, (usr.FirstName || ' ' || usr.LastName) AS departmentMember, 
		un.CreatedAt, un.Remarks, un.NoticeEndDate, un.IsApproved, 
		(apusr.FirstName || ' ' || apusr.LastName) AS approvedBy
		FROM UserNotice un
		INNER JOIN DepartmentMember dm ON dm.ID = un.departmentMemberID AND dm.IsActive = 1 
		INNER JOIN [User] usr ON dm.UserID = usr.ID AND usr.IsActive = 1
		LEFT JOIN [User] apusr ON un.approvedBy = apusr.ID AND apusr.IsActive = 1
		WHERE NoticeEndDate > ?
		LIMIT ? OFFSET ?`, now, itemsPerPage, (filters.Page-1)*uint(itemsPerPage)).
		Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *noticeRepository) GetNoticeUserCount() (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM UserNotice
		WHERE IsActive = 1 AND NoticeEndDate > ?`, time.Now()).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *noticeRepository) GetNoticeUserCountByDepartment(departmentID uint) (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM UserNotice
		INNER JOIN DepartmentMember dm ON dm.ID = un.departmentMemberID AND dm.IsActive = 1 
		WHERE IsActive = 1 AND NoticeEndDate > ? AND dm.DepartmentID = ?`,
		time.Now(), departmentID).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *noticeRepository) FetchNotice(departmentMemberID uint) (*response.FetchActiveUserNotices, error) {
	var data response.FetchActiveUserNotices

	if err := r.db.Raw(`
		SELECT un.DepartmentMemberID, (usr.FirstName || ' ' || usr.LastName) AS departmentMember, 
		un.CreatedAt, un.Remarks, un.NoticeEndDate, IsApproved, 
		(apusr.FirstName || ' ' || apusr.LastName) AS approvedBy
		FROM UserNotice un
		INNER JOIN DepartmentMember dm ON dm.ID = un.departmentMemberID AND dm.IsActive = 1 
		INNER JOIN [User] usr ON dm.UserID = usr.ID AND usr.IsActive = 1
		LEFT JOIN [User] apusr ON un.approvedBy = apusr.ID AND apusr.IsActive = 1
		WHERE un.IsActive = 1 AND un.DepartmentMemberID = ?`, departmentMemberID).
		Scan(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *noticeRepository) ApproveNotice(departmentMemberID, approvedBy uint, req *request.ApproveNotice) error {
	var appliedDate time.Time

	if err := r.db.Raw(`
		SELECT CreatedAt
		FROM UserNotice
		WHERE DepartmentMemberID = ?`, departmentMemberID).Scan(&appliedDate).Error; err != nil {
		return err
	}

	noticeEndDate := appliedDate.AddDate(0, 0, req.ServeDays)

	return r.db.Exec(`
		UPDATE UserNotice
		SET UpdatedAt = ?, NoticeEndDate = ?, IsApproved = 1, ApprovedBy = ?
		WHERE DepartmentMemberID = ?`,
		time.Now(), noticeEndDate, approvedBy, departmentMemberID).Error
}
