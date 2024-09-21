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

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domain.PermissionRepository {
	return &permissionRepository{db}
}

func (r *permissionRepository) RequestPermission(departmentMemberID uint, req *request.RequestPermission) error {
	return r.db.Exec(`
			INSERT INTO DepartmentMemberPermissionRequest
			(CreatedAt, UpdatedAt, DepartmentMemberID, [Date], FromTime, ToTime, Reason)
			VALUES(?, ?, ?, ?, ?, ?, ?)`, time.Now(), time.Now(),
		departmentMemberID, req.Date, req.FromTime, req.ToTime, req.Reason).Error
}

func (r *permissionRepository) FetchPermissions(departmentMemberID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	var data []response.FetchUserPermissions
	var itemsPerPage uint = 10
	var totalCount int = 0
	var baseQuery strings.Builder
	var queryParams []interface{}

	baseQuery.WriteString(`
		SELECT dmpr.ID, dmpr.DepartmentMemberID departmentMemberID, dmpr.FromTime fromTime, 
		dmpr.ToTime toTime, dmpr.Reason, dmpr.IsApproved, dmpr.ApprovedAt, dmpr.CreatedAt, dmpr.date,
		dmpr.UpdatedAt, dmpr.IsActive, COUNT(1) OVER (PARTITION BY 1) AS [count],
		(deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember, 
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy
		FROM DepartmentMemberPermissionRequest dmpr
		INNER JOIN DepartmentMember dm ON dm.ID = dmpr.departmentMemberID AND dm.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmpr.ApprovedBy AND approvedUser.IsActive
		WHERE dmpr.IsActive = 1 AND dmpr.DepartmentMemberID = ?
		ORDER BY dmpr.CreatedAt DESC`)

	queryParams = append(queryParams, departmentMemberID)

	if filters.Page > 0 {
		baseQuery.WriteString(` LIMIT ? OFFSET ?`)
		queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*itemsPerPage)
	}

	if err := r.db.Raw(baseQuery.String(), queryParams...).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count
		permissionCount, err := r.GetPermissionCountByUser(departmentMemberID)

		if err != nil {
			return nil, err
		}

		data[0].UserPermissionCount = permissionCount
	}

	response := *utils.PaginatedResponse(uint(totalCount), filters.Page, data)

	return &response, nil
}

func (r *permissionRepository) FetchUserPermissions(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	var data []response.FetchUserPermissions
	var itemsPerPage uint = 10
	var totalCount int = 0
	var baseQuery strings.Builder
	var queryParams []interface{}

	baseQuery.WriteString(`
		SELECT dmpr.ID, dmpr.DepartmentMemberID departmentMemberID, dmpr.ApprovedAt, 
		dmpr.FromTime fromTime, dmpr.ToTime toTime, dmpr.Reason, dmpr.IsActive, dmpr.CreatedAt,
		dmpr.UpdatedAt, dmpr.IsApproved, COUNT(1) OVER (PARTITION BY 1) AS [count], dmpr.date,
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy,
		(deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember
		FROM DepartmentMemberPermissionRequest dmpr
		INNER JOIN DepartmentMember dm ON dm.ID = dmpr.departmentMemberID AND dm.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmpr.ApprovedBy AND approvedUser.IsActive
		WHERE dmpr.IsActive = 1 AND (dm.DepartmentID = ? OR dm.departmentID = 0)
		ORDER BY dmpr.CreatedAt DESC`)

	queryParams = append(queryParams, departmentID)

	if filters.Page > 0 {
		baseQuery.WriteString(` LIMIT ? OFFSET ?`)
		queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*itemsPerPage)
	}

	if err := r.db.Raw(baseQuery.String(), queryParams...).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count
		data[0].UserPermissionCount = totalCount
	}

	response := *utils.PaginatedResponse(uint(totalCount), filters.Page, data)

	return &response, nil
}

func (r *permissionRepository) UpdatePermissionStatus(permissionID, approvedBy uint, req *request.UpdatePermissionStatus) error {
	return r.db.Exec(`
		UPDATE DepartmentMemberPermissionRequest
		SET UpdatedAt = ?, IsApproved = ?, ApprovedAt = ?, ApprovedBy = ?
		WHERE ID = ?`, time.Now(), req.IsApproved, time.Now(), approvedBy, permissionID).Error
}

func (r *permissionRepository) IsPermissionExistWithoutApproval(departmentMemberID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberPermissionRequest
		WHERE DepartmentMemberID = ? AND IsApproved IS NULL AND IsActive = 1
		AND date(CreatedAt) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(CreatedAt) <= date(strftime('%Y-%m', 'now') || '-26');`, departmentMemberID).
		Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *permissionRepository) IsPermissionExistWithID(id uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberPermissionRequest
		WHERE ID = ? AND IsActive = 1`, id).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *permissionRepository) UpdatePermissionRequest(permissionID uint, req *request.RequestPermission) error {
	return r.db.Exec(`
		UPDATE DepartmentMemberPermissionRequest
		SET UpdatedAt = ?, Reason = ?, [Date] = ?, FromTime = ?, ToTime = ?
		WHERE ID = ?`, time.Now(), req.Reason, req.Date, req.FromTime, req.ToTime, permissionID).Error
}

func (r *permissionRepository) RemovePermissionRequest(permissionID uint) error {
	return r.db.Exec(`
		UPDATE DepartmentMemberPermissionRequest
		SET IsActive = ?, DeletedAt = ?
		WHERE ID = ?`, constant.Inactive, time.Now(), permissionID).Error
}

func (r *permissionRepository) IsPermissionExistsWithApproval(permissionID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberPermissionRequest dmpr
		WHERE dmpr.ID = ? AND dmpr.IsApproved = 1`, permissionID).
		Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *permissionRepository) GetPermissionCount() (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberPermissionRequest
		WHERE IsActive = 1
		AND date(CreatedAt) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(CreatedAt) <= date(strftime('%Y-%m', 'now') || '-26')`).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *permissionRepository) GetPermissionCountByUser(departmentMemberID uint) (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberPermissionRequest
		WHERE IsActive = 1 AND DeparmentMemberID = ?
		AND date(CreatedAt) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(CreatedAt) <= date(strftime('%Y-%m', 'now') || '-26')`, departmentMemberID).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *permissionRepository) GetPermissionCountByDepartment(departmentID uint) (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberPermissionRequest dmpr
		INNER JOIN DepartmentMember dm ON dm.id = dmpr.DepartmentMemberID AND dm.IsActive = 1
		WHERE dmpr.IsActive = 1 AND dm.DepartmentID = ? AND dmpr.IsApproved = 1
		AND date(CreatedAt) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(CreatedAt) <= date(strftime('%Y-%m', 'now') || '-26')`, departmentID).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *permissionRepository) GetApprovedPermissionCount() (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberPermissionRequest
		WHERE IsActive = 1 AND IsApproved = 1
		AND date(CreatedAt) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(CreatedAt) <= date(strftime('%Y-%m', 'now') || '-26')`).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *permissionRepository) GetApprovedPermissionCountByDepartment(departmentID uint) (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberPermissionRequest dmpr
		INNER JOIN DepartmentMember dm ON dm.id = dmpr.DepartmentMemberID AND dm.IsActive = 1
		WHERE dmpr.IsActive = 1 AND dm.DepartmentID = ?
		AND date(CreatedAt) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(CreatedAt) <= date(strftime('%Y-%m', 'now') || '-26')`, departmentID).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *permissionRepository) GetApprovedPermissionCountByUser(departmentMemberID uint) (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberPermissionRequest
		WHERE IsActive = 1 AND IsApproved = 1 AND DeparmentMemberID = ?
		AND date(CreatedAt) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(CreatedAt) <= date(strftime('%Y-%m', 'now') || '-26')`, departmentMemberID).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
