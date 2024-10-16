package repository

import (
	"ems/app/model/constant"
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/domain"
	"ems/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type noticeRepository struct {
	db *gorm.DB
}

func NewNoticeRepository(db *gorm.DB) domain.NoticeRepository {
	return &noticeRepository{db}
}

func (r *noticeRepository) ApplyNotice(departmentMemberID uint, req *request.ApplyNotice) error {
	return r.db.Exec(`
			INSERT INTO UserNotice
			(CreatedAt, UpdatedAt, DepartmentMemberID, Remarks)
			VALUES(?, ?, ?, ?)`,
		time.Now(), time.Now(), departmentMemberID, req.Remarks).Error
}

func (r *noticeRepository) FetchActiveUserNotices(roleID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	var (
		data         []response.FetchActiveUserNotices
		itemsPerPage uint = 10
		totalCount   int  = 0
		query        strings.Builder
		queryParams  []interface{}
	)

	query.WriteString(`
		SELECT un.DepartmentMemberID, (usr.FirstName || ' ' || usr.LastName) AS departmentMember, 
		un.CreatedAt, un.Remarks, un.NoticeEndDate, un.IsApproved, (apusr.FirstName || ' ' || apusr.LastName) AS approvedBy,
		COUNT(*) OVER (PARTITION BY 1) AS [count], [Role].[Name] AS [role]
		FROM UserNotice un
		INNER JOIN DepartmentMember dm ON dm.ID = un.departmentMemberID AND dm.IsActive = 1 
		INNER JOIN [User] usr ON dm.UserID = usr.ID AND usr.IsActive = 1
		INNER JOIN [Role] ON [Role].ID = usr.RoleID AND [Role].IsActive = 1
		LEFT JOIN [User] apusr ON un.approvedBy = apusr.ID AND apusr.IsActive = 1`)

	if roleID == uint(constant.HR) {
		query.WriteString(` WHERE usr.RoleID = ?`)
		queryParams = append(queryParams, constant.Employee)
	} else {
		query.WriteString(` WHERE usr.RoleID IN ?`)
		queryParams = append(queryParams, []interface{}{constant.HR, constant.DepartmentLead})
	}

	query.WriteString(` ORDER BY un.CreatedAt DESC LIMIT ? OFFSET ?`)
	queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*uint(itemsPerPage))

	if err := r.db.Raw(query.String(), queryParams...).
		Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count
	}

	response := *utils.PaginatedResponse(uint(totalCount), filters.Page, data)

	return &response, nil
}

func (r *noticeRepository) GetNoticeUserCount() (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(*)
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
		SELECT COUNT(*)
		FROM UserNotice un
		INNER JOIN DepartmentMember dm ON dm.ID = un.departmentMemberID AND dm.IsActive = 1 
		WHERE un.IsActive = 1 AND un.NoticeEndDate > ? AND dm.DepartmentID = ?`,
		time.Now(), departmentID).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *noticeRepository) FetchNotice(departmentMemberID uint) (*response.FetchActiveUserNotices, error) {
	var data *response.FetchActiveUserNotices

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

	return data, nil
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

func (r *noticeRepository) IsApproveExistsByUser(departmentMemberID uint) (bool, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM UserNotice
		WHERE DepartmentMemberID = ?`, departmentMemberID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
