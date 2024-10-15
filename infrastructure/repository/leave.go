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

type leaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) domain.LeaveRepository {
	return &leaveRepository{db}
}

func (r *leaveRepository) RequestLeave(departmentMemberID uint, req *request.RequestLeave) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			INSERT INTO DepartmentMemberLeaveRequest
			(CreatedAt, UpdatedAt, DepartmentMemberID, Reason)
			VALUES(?, ?, ?, ?)`, time.Now(), time.Now(), departmentMemberID, req.Reason).Error; err != nil {
			return err
		}

		var departmentMemberLeaveRequestID uint

		if err := tx.Raw(`
			SELECT ID
			FROM DepartmentMemberLeaveRequest
			ORDER BY CreatedAt DESC LIMIT 1`).Scan(&departmentMemberLeaveRequestID).Error; err != nil {
			return err
		}

		for _, date := range req.Dates {
			if err := tx.Exec(`
			INSERT INTO DepartmentMemberLeaveRequestDate
			(CreatedAt, UpdatedAt, DepartmentMemberLeaveRequestID, [date], IsFullDay, SessionType)
			VALUES(?, ?, ?, ?, ?, ?)`, time.Now(), time.Now(), departmentMemberLeaveRequestID,
				date.Date, date.IsFullDay, date.SessionType).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *leaveRepository) FetchOwnLeaves(departmentMemberID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	var (
		data         []response.FetchLeaves
		itemsPerPage uint = 10
		totalCount   uint = 0
		query        strings.Builder
		queryParams  []interface{}
	)
	startDate, endDate := utils.GetDateRangeForMonthAndYear(filters.Year, filters.Month)

	query.WriteString(`
		SELECT dmlr.ID, dmlr.DepartmentMemberID departmentMemberID, dmlr.ApprovedAt, 
		dmlr.Reason, GROUP_CONCAT(strftime('%Y-%m-%d', dmlrd.[Date])) AS dates, dmlr.UpdatedAt, dmlr.IsActive, 
		GROUP_CONCAT(dmlrd.IsFullDay) AS isFullDays, GROUP_CONCAT(dmlrd.SessionType) AS sessionTypes, 
		COUNT(*) OVER (PARTITION BY 1) AS [count], dmlr.CreatedAt, dmlr.IsApproved, (deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember,
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		INNER JOIN DepartmentMember dm ON dm.ID = dmlr.departmentMemberID AND dm.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmlr.ApprovedBy AND approvedUser.IsActive
		WHERE dmlr.IsActive = 1 AND dmlr.DepartmentMemberID = ? AND dmlrd.Date BETWEEN ? AND ?
		GROUP BY dmlr.ID, dmlr.DepartmentMemberID, dmlr.Reason, dmlr.CreatedAt, 
		dmlr.UpdatedAt, dmlr.IsActive
		ORDER BY dmlr.CreatedAt DESC`)

	queryParams = append(queryParams, departmentMemberID, startDate, endDate)

	if filters.Page > 0 {
		query.WriteString(` LIMIT ? OFFSET ?`)
		queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*itemsPerPage)
	}

	if err := r.db.Raw(query.String(), queryParams...).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count
		leaveCount, err := r.GetLeaveCountByUser(departmentMemberID, &filters.DateFilters)
		if err != nil {
			return nil, err
		}
		data[0].UserLeaveCount = leaveCount
	}

	response := *utils.PaginatedResponse(totalCount, filters.Page, data)
	return &response, nil
}

func (r *leaveRepository) FetchDepartmentMemberLeaves(departmentID uint, filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	var (
		data         []response.FetchLeaves
		itemsPerPage uint = 10
		totalCount   uint = 0
	)
	startDate, endDate := utils.GetDateRangeForMonthAndYear(filters.Year, filters.Month)

	if err := r.db.Raw(`
		SELECT dmlr.ID, dmlr.DepartmentMemberID departmentMemberID, dmlr.ApprovedAt, 
		dmlr.Reason, GROUP_CONCAT(strftime('%Y-%m-%d', dmlrd.[Date])) AS dates, dmlr.UpdatedAt, dmlr.IsActive, 
		GROUP_CONCAT(dmlrd.IsFullDay) AS isFullDays, GROUP_CONCAT(dmlrd.SessionType) AS sessionTypes, 
		COUNT(*) OVER (PARTITION BY 1) AS [count],dmlr.CreatedAt, dmlr.IsApproved, 
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy,
		(deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember, [Role].[Name] AS [role]
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMember dm ON dm.ID = dmlr.departmentMemberID AND dm.IsActive = 1 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		INNER JOIN [Role] ON [Role].ID = deptMem.RoleID AND [Role].IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmlr.ApprovedBy AND approvedUser.IsActive
		WHERE dmlr.IsActive = 1 AND dm.DepartmentID = ? AND deptMem.RoleID NOT IN ? 
		AND dmlrd.Date BETWEEN ? AND ?
		GROUP BY dmlr.ID, dmlr.DepartmentMemberID, dmlr.Reason, dmlr.CreatedAt, 
		dmlr.UpdatedAt, dmlr.IsActive 
		ORDER BY dmlr.CreatedAt DESC 
		LIMIT ? OFFSET ?`, departmentID, []interface{}{constant.DepartmentLead, constant.HR},
		startDate, endDate, itemsPerPage, (filters.Page-1)*itemsPerPage).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count

		for i := range data {
			leaveCount, err := r.GetLeaveCountByUser(data[i].DepartmentMemberID, &filters.DateFilters)
			if err != nil {
				return nil, err
			}
			data[i].UserLeaveCount = leaveCount
		}

	}

	response := *utils.PaginatedResponse(totalCount, filters.Page, data)

	return &response, nil
}

func (r *leaveRepository) UpdateLeaveStatus(leaveID, approvedBy uint, req *request.UpdateLeaveStatus) error {
	return r.db.Exec(`
		UPDATE DepartmentMemberLeaveRequest 
		SET UpdatedAt = ?, IsApproved = ?, ApprovedAt = ?, ApprovedBy = ?
		WHERE ID = ?`, time.Now(), req.IsApproved, time.Now(), approvedBy, leaveID).Error
}

func (r *leaveRepository) IsLeaveExistsWithoutApproval(departmentMemberID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.DepartmentMemberID = ? AND dmlr.IsApproved IS NULL AND dmlr.IsActive = 1`,
		departmentMemberID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *leaveRepository) FetchLeadAndHRLeaves(filters *request.CommonRequestWithDateFilter) (*utils.PaginationResponse, error) {
	var (
		data         []response.FetchLeaves
		itemsPerPage uint = 10
		totalCount   uint = 0
	)

	startDate, endDate := utils.GetDateRangeForMonthAndYear(filters.Year, filters.Month)

	if err := r.db.Raw(`
		SELECT dmlr.ID, dmlr.DepartmentMemberID departmentMemberID, dmlr.ApprovedAt, 
		dmlr.Reason, GROUP_CONCAT(strftime('%Y-%m-%d', dmlrd.[Date])) AS dates, dmlr.UpdatedAt, 
		dmlr.IsActive, GROUP_CONCAT(dmlrd.IsFullDay) AS isFullDays, GROUP_CONCAT(dmlrd.SessionType) AS sessionTypes, 
		COUNT(*) OVER (PARTITION BY 1) AS [count], dmlr.CreatedAt, dmlr.IsApproved,  [Role].[Name] AS [role],
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy,
		(deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMember dm ON dm.ID = dmlr.departmentMemberID AND dm.IsActive = 1 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		INNER JOIN [Role] ON [Role].ID = deptMem.RoleID AND [Role].IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmlr.ApprovedBy AND approvedUser.IsActive
		WHERE dmlr.IsActive = 1 AND deptMem.RoleID IN ? AND dmlrd.Date BETWEEN ? AND ?
		GROUP BY dmlr.ID, dmlr.DepartmentMemberID, dmlr.Reason, dmlr.CreatedAt, 
		dmlr.UpdatedAt, dmlr.IsActive
		ORDER BY dmlr.CreatedAt DESC
		LIMIT ? OFFSET ?`, []interface{}{constant.DepartmentLead, constant.HR},
		startDate, endDate, itemsPerPage, (filters.Page-1)*itemsPerPage).
		Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count

		for i := range data {
			leaveCount, err := r.GetLeaveCountByUser(data[i].DepartmentMemberID,
				&filters.DateFilters)
			if err != nil {
				return nil, err
			}
			data[i].UserLeaveCount = leaveCount
		}
	}

	response := *utils.PaginatedResponse(totalCount, filters.Page, data)

	return &response, nil
}

func (r *leaveRepository) UpdateLeaveRequest(leaveID uint, req *request.RequestLeave) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			UPDATE DepartmentMemberLeaveRequest
			SET UpdatedAt = ?, Reason = ?
			WHERE ID = ?`, time.Now(), req.Reason, leaveID).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE DepartmentMemberLeaveRequestDate
			SET UpdatedAt = ?, IsActive = ?
			WHERE DepartmentMemberLeaveRequestID = ?`, time.Now(), constant.Inactive, leaveID).
			Error; err != nil {
			return err
		}

		for _, date := range req.Dates {
			if err := tx.Exec(`
			INSERT INTO DepartmentMemberLeaveRequestDate
			(CreatedAt, UpdatedAt, DepartmentMemberLeaveRequestID, [date], IsFullDay, SessionType)
			VALUES(?, ?, ?, ?, ?, ?)`, time.Now(), time.Now(), leaveID,
				date.Date, date.IsFullDay, date.SessionType).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *leaveRepository) IsLeaveExistsWithID(leaveID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM DepartmentMemberLeaveRequest
		WHERE ID = ? AND IsActive = 1`, leaveID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *leaveRepository) RemoveLeaveRequest(leaveID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			UPDATE DepartmentMemberLeaveRequest
			SET IsActive = ?, DeletedAt = ?
			WHERE ID = ?`, constant.Inactive, time.Now(), leaveID).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE DepartmentMemberLeaveRequestDate
			SET IsActive = ?, DeletedAt = ?
			WHERE DepartmentMemberLeaveRequestID = ?`, constant.Inactive, time.Now(), leaveID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *leaveRepository) IsLeaveExistsWithApproval(leaveID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.ID = ? AND dmlr.IsApproved = 1`, leaveID).
		Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *leaveRepository) GetLeaveCountByUser(departmentMemberID uint, dateFilters *request.DateFilters) (float64, error) {
	var (
		data []struct {
			ID        uint `gorm:"column:ID"`
			IsFullDay bool `gorm:"column:IsFullDay"`
		}
		leaveCount float64
	)
	startDate, endDate := utils.GetDateRangeForMonthAndYear(dateFilters.Year, dateFilters.Month)

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.DepartmentMemberID = ? AND dmlr.IsActive = 1 AND dmlrd.[Date] BETWEEN ? AND ?`,
		departmentMemberID, startDate, endDate).
		Scan(&data).Error; err != nil {
		return 0, err
	}

	for _, leave := range data {
		if leave.IsFullDay {
			leaveCount += 1
			continue
		}
		leaveCount += 0.5
	}

	return leaveCount, nil
}

func (r *leaveRepository) GetLeaveCount(dateFilters *request.DateFilters) (float64, error) {
	var (
		data []struct {
			ID        uint `gorm:"column:ID"`
			IsFullDay bool `gorm:"column:IsFullDay"`
		}
		leaveCount float64
	)

	startDate, endDate := utils.GetDateRangeForMonthAndYear(dateFilters.Year, dateFilters.Month)

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.IsActive = 1 AND dmlrd.[Date] BETWEEN ? AND ?`, startDate, endDate).
		Scan(&data).Error; err != nil {
		return 0, err
	}

	for _, leave := range data {
		if leave.IsFullDay {
			leaveCount += 1
			continue
		}
		leaveCount += 0.5
	}

	return leaveCount, nil
}

func (r *leaveRepository) GetApprovedLeaveCount(dateFilters *request.DateFilters) (float64, error) {
	var (
		data []struct {
			ID        uint `gorm:"column:ID"`
			IsFullDay bool `gorm:"column:IsFullDay"`
		}
		leaveCount float64
	)

	startDate, endDate := utils.GetDateRangeForMonthAndYear(dateFilters.Year, dateFilters.Month)

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.IsActive = 1 AND dmlr.IsApproved = 1 AND dmlrd.[Date] BETWEEN ? AND ?`,
		startDate, endDate).
		Scan(&data).Error; err != nil {
		return 0, err
	}

	for _, leave := range data {
		if leave.IsFullDay {
			leaveCount += 1
			continue
		}
		leaveCount += 0.5
	}

	return leaveCount, nil
}

func (r *leaveRepository) GetApprovedLeaveCountByUser(departmentMemberID uint, dateFilters *request.DateFilters) (float64, error) {
	var (
		data []struct {
			ID        uint `gorm:"column:ID"`
			IsFullDay bool `gorm:"column:IsFullDay"`
		}
		leaveCount float64
	)
	startDate, endDate := utils.GetDateRangeForMonthAndYear(dateFilters.Year, dateFilters.Month)

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.IsActive = 1 AND dmlr.IsApproved = 1 AND dmlr.DepartmentMemberID = ?
		AND dmlrd.[Date] BETWEEN ? AND ?`, departmentMemberID, startDate, endDate).
		Scan(&data).Error; err != nil {
		return 0, err
	}

	for _, leave := range data {
		if leave.IsFullDay {
			leaveCount += 1
			continue
		}
		leaveCount += 0.5
	}

	return leaveCount, nil
}
