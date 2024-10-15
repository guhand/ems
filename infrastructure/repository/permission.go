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

func (r *permissionRepository) FetchOwnPermissions(departmentMemberID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	var (
		data         []response.FetchUserPermissions
		itemsPerPage uint = 10
		totalCount   int  = 0
		query        strings.Builder
		queryParams  []interface{}
		now          time.Time = time.Now()
	)
	startDate, endDate := utils.GetDateRangeForMonthAndYear(filters.Year, filters.Month)

	query.WriteString(`
		SELECT dmpr.ID, dmpr.DepartmentMemberID departmentMemberID, dmpr.FromTime fromTime, 
		dmpr.ToTime toTime, dmpr.Reason, dmpr.IsApproved, dmpr.ApprovedAt, dmpr.CreatedAt, 
		dmpr.UpdatedAt, dmpr.IsActive, COUNT(*) OVER (PARTITION BY 1) AS [count], 
		strftime('%Y-%m-%d', dmpr.[Date]) AS [date], (deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember, 
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy
		FROM DepartmentMemberPermissionRequest dmpr
		INNER JOIN DepartmentMember dm ON dm.ID = dmpr.departmentMemberID AND dm.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmpr.ApprovedBy AND approvedUser.IsActive
		WHERE dmpr.IsActive = 1 AND dmpr.DepartmentMemberID = ? AND dmpr.[Date] BETWEEN ? AND ?`)

	queryParams = append(queryParams, departmentMemberID, startDate, endDate)

	if filters.Page > 0 {
		query.WriteString(` ORDER BY dmpr.CreatedAt DESC LIMIT ? OFFSET ?`)
		queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*itemsPerPage)
	}

	if err := r.db.Raw(query.String(), queryParams...).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count
		filters.DateFilters.Year = now.Year()
		filters.DateFilters.Month = int(now.Month())

		if now.Day() > 26 {
			if now.Month() == 12 {
				filters.DateFilters.Month = 1
				filters.Year += 1
			} else {
				filters.DateFilters.Month += 1
			}
		}
		permissionCount, err := r.GetPermissionCountByUser(departmentMemberID, &filters.DateFilters)

		if err != nil {
			return nil, err
		}

		data[0].UserPermissionCount = permissionCount
	}

	response := *utils.PaginatedResponse(uint(totalCount), filters.Page, data)

	return &response, nil
}

func (r *permissionRepository) FetchDepartmentMemberPermissions(departmentID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	var (
		data         []response.FetchUserPermissions
		itemsPerPage uint = 10
		totalCount   int  = 0
	)

	startDate, endDate := utils.GetDateRangeForMonthAndYear(filters.Year, filters.Month)

	if err := r.db.Raw(`
		SELECT dmpr.ID, dmpr.DepartmentMemberID departmentMemberID, dmpr.ApprovedAt, 
		dmpr.FromTime fromTime, dmpr.ToTime toTime, dmpr.Reason, dmpr.IsActive, dmpr.CreatedAt,
		dmpr.UpdatedAt, dmpr.IsApproved, COUNT(*) OVER (PARTITION BY 1) AS [count], 
		strftime('%Y-%m-%d', dmpr.[Date]) AS [date], (approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy,
		(deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember, [Role].[Name] AS [role]
		FROM DepartmentMemberPermissionRequest dmpr
		INNER JOIN DepartmentMember dm ON dm.ID = dmpr.departmentMemberID AND dm.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		INNER JOIN [Role] ON [Role].ID = deptMem.RoleID AND [Role].IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmpr.ApprovedBy AND approvedUser.IsActive
		WHERE dmpr.IsActive = 1 AND dm.DepartmentID = ? AND deptMem.RoleID <> ? 
		AND dmpr.[Date] BETWEEN ? AND ?
		ORDER BY dmpr.CreatedAt DESC LIMIT ? OFFSET ?`, departmentID, constant.DepartmentLead, startDate, endDate,
		itemsPerPage, (filters.Page-1)*itemsPerPage).Scan(&data).Error; err != nil {
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
		SELECT COUNT(*)
		FROM DepartmentMemberPermissionRequest
		WHERE DepartmentMemberID = ? AND IsApproved IS NULL AND IsActive = 1`,
		departmentMemberID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *permissionRepository) IsPermissionExistWithID(id uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
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
		SELECT COUNT(*)
		FROM DepartmentMemberPermissionRequest
		WHERE ID = ? AND IsApproved = 1`, permissionID).
		Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *permissionRepository) GetPermissionCount(dateFilters *request.DateFilters) (int, error) {
	var count int
	startDate, endDate := utils.GetDateRangeForMonthAndYear(dateFilters.Year, dateFilters.Month)

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM DepartmentMemberPermissionRequest
		WHERE IsActive = 1 AND [Date] BETWEEN ? AND ?`, startDate, endDate).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *permissionRepository) GetPermissionCountByUser(departmentMemberID uint, dateFilters *request.DateFilters) (int, error) {
	var count int
	startDate, endDate := utils.GetDateRangeForMonthAndYear(dateFilters.Year, dateFilters.Month)

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM DepartmentMemberPermissionRequest
		WHERE IsActive = 1 AND DepartmentMemberID = ? AND [Date] BETWEEN ? AND ?`,
		departmentMemberID, startDate, endDate).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *permissionRepository) GetApprovedPermissionCount(dateFilters *request.DateFilters) (int, error) {
	var count int
	startDate, endDate := utils.GetDateRangeForMonthAndYear(dateFilters.Year, dateFilters.Month)

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM DepartmentMemberPermissionRequest
		WHERE IsActive = 1 AND IsApproved = 1 AND [Date] BETWEEN ? AND ?`, startDate, endDate).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *permissionRepository) GetApprovedPermissionCountByUser(departmentMemberID uint, dateFilters *request.DateFilters) (int, error) {
	var count int
	startDate, endDate := utils.GetDateRangeForMonthAndYear(dateFilters.Year, dateFilters.Month)

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM DepartmentMemberPermissionRequest
		WHERE IsActive = 1 AND IsApproved = 1 AND DepartmentMemberID = ? 
		AND [Date] BETWEEN ? AND ?`, departmentMemberID, startDate, endDate).
		Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *permissionRepository) FetchLeadAndHRPermissions(filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	var (
		data         []response.FetchUserPermissions
		itemsPerPage uint = 10
		totalCount   int  = 0
	)
	startDate, endDate := utils.GetDateRangeForMonthAndYear(filters.Year, filters.Month)

	if err := r.db.Raw(`
		SELECT dmpr.ID, dmpr.DepartmentMemberID departmentMemberID, dmpr.ApprovedAt, 
		dmpr.FromTime fromTime, dmpr.ToTime toTime, dmpr.Reason, dmpr.IsActive, dmpr.CreatedAt,
		dmpr.UpdatedAt, dmpr.IsApproved, COUNT(*) OVER (PARTITION BY 1) AS [count],
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy, [Role].[Name] AS [role],
		strftime('%Y-%m-%d', dmpr.[Date]) AS [date], (deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember
		FROM DepartmentMemberPermissionRequest dmpr
		INNER JOIN DepartmentMember dm ON dm.ID = dmpr.departmentMemberID AND dm.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		INNER JOIN [Role] ON [Role].ID = deptMem.RoleID AND [Role].IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmpr.ApprovedBy AND approvedUser.IsActive
		WHERE dmpr.IsActive = 1 AND deptMem.RoleID IN ? AND dmpr.Date BETWEEN ? AND ?
		ORDER BY dmpr.CreatedAt DESC LIMIT ? OFFSET ?`,
		[]interface{}{constant.DepartmentLead, constant.HR}, startDate, endDate,
		itemsPerPage, (filters.Page-1)*itemsPerPage).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count
		data[0].UserPermissionCount = totalCount
	}

	response := *utils.PaginatedResponse(uint(totalCount), filters.Page, data)

	return &response, nil
}
