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

func (r *leaveRepository) FetchOwnLeaves(departmentMemberID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	var data []response.FetchLeaves
	var itemsPerPage uint = 10
	var totalCount uint = 0
	var baseQuery strings.Builder
	var queryParams []interface{}

	baseQuery.WriteString(`
		SELECT dmlr.ID, dmlr.DepartmentMemberID departmentMemberID, dmlr.ApprovedAt, 
		dmlr.Reason, GROUP_CONCAT(dmlrd.Date) AS dates, dmlr.UpdatedAt, dmlr.IsActive, 
		GROUP_CONCAT(dmlrd.IsFullDay) AS isFullDays, GROUP_CONCAT(dmlrd.SessionType) AS sessionTypes, 
		COUNT(1) OVER (PARTITION BY 1) AS [count], dmlr.CreatedAt, dmlr.IsApproved, 
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmlr.ApprovedBy AND approvedUser.IsActive
		WHERE dmlr.IsActive = 1 AND dmlr.DepartmentMemberID = ?
		GROUP BY dmlr.ID, dmlr.DepartmentMemberID, dmlr.Reason, dmlr.CreatedAt, 
		dmlr.UpdatedAt, dmlr.IsActive
		ORDER BY dmlr.CreatedAt DESC`)

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
		leaveCount, err := r.GetLeaveCountByUser(departmentMemberID)

		if err != nil {
			return nil, err
		}

		data[0].UserLeaveCount = leaveCount
	}

	response := *utils.PaginatedResponse(totalCount, filters.Page, data)

	return &response, nil
}

func (r *leaveRepository) FetchDepartmentMemberLeaves(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	var data []response.FetchLeaves
	var itemsPerPage uint = 10
	var totalCount uint = 0
	var baseQuery strings.Builder
	var queryParams []interface{}

	baseQuery.WriteString(`
		SELECT dmlr.ID, dmlr.DepartmentMemberID departmentMemberID, dmlr.ApprovedAt, 
		dmlr.Reason, GROUP_CONCAT(dmlrd.Date) AS dates, dmlr.UpdatedAt, dmlr.IsActive, 
		GROUP_CONCAT(dmlrd.IsFullDay) AS isFullDays, GROUP_CONCAT(dmlrd.SessionType) AS sessionTypes, 
		COUNT(1) OVER (PARTITION BY 1) AS [count],dmlr.CreatedAt, dmlr.IsApproved, 
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy,
		(deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMember dm ON dm.ID = dmlr.departmentMemberID AND dm.IsActive = 1 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmlr.ApprovedBy AND approvedUser.IsActive
		WHERE dmlr.IsActive = 1 AND dm.DepartmentID = ? AND deptMem.RoleID NOT IN ?
		GROUP BY dmlr.ID, dmlr.DepartmentMemberID, dmlr.Reason, dmlr.CreatedAt, 
		dmlr.UpdatedAt, dmlr.IsActive
		ORDER BY dmlr.CreatedAt DESC`)

	queryParams = append(queryParams, departmentID, []interface{}{constant.DepartmentLead, constant.HR})

	if filters.Page > 0 {
		baseQuery.WriteString(` LIMIT ? OFFSET ?`)
		queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*itemsPerPage)
	}

	if err := r.db.Raw(baseQuery.String(), queryParams...).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count

		for i := range data {
			leaveCount, err := r.GetLeaveCountByUser(data[i].DepartmentMemberID)
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
		SELECT COUNT(1)
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.DepartmentMemberID = ? AND dmlr.IsApproved IS NULL AND IsActive = 1
		AND date(dmlrd.[Date]) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(dmlrd.[Date]) <= date(strftime('%Y-%m', 'now') || '-26');`, departmentMemberID).
		Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *leaveRepository) FetchLeaves(filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	var data []response.FetchLeaves
	var itemsPerPage uint = 10
	var totalCount uint = 0
	var baseQuery strings.Builder
	var queryParams []interface{}

	baseQuery.WriteString(`
		SELECT dmlr.ID, dmlr.DepartmentMemberID departmentMemberID, dmlr.ApprovedAt, 
		dmlr.Reason, GROUP_CONCAT(dmlrd.Date) AS dates, dmlr.UpdatedAt, dmlr.IsActive, 
		GROUP_CONCAT(dmlrd.IsFullDay) AS isFullDays, GROUP_CONCAT(dmlrd.SessionType) AS sessionTypes, 
		COUNT(1) OVER (PARTITION BY 1) AS [count],dmlr.CreatedAt, dmlr.IsApproved, 
		(approvedUser.FirstName || ' ' || approvedUser.LastName) AS approvedBy,
		(deptMem.FirstName || ' ' || deptMem.LastName) AS departmentMember
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMember dm ON dm.ID = dmlr.departmentMemberID AND dm.IsActive = 1 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		INNER JOIN [User] deptMem ON dm.UserID = deptMem.ID AND deptMem.IsActive = 1
		LEFT JOIN [User] approvedUser ON approvedUser.ID = dmlr.ApprovedBy AND approvedUser.IsActive
		WHERE dmlr.IsActive = 1
		GROUP BY dmlr.ID, dmlr.DepartmentMemberID, dmlr.Reason, dmlr.CreatedAt, 
		dmlr.UpdatedAt, dmlr.IsActive
		ORDER BY dmlr.CreatedAt DESC`)

	if filters.Page > 0 {
		baseQuery.WriteString(` LIMIT ? OFFSET ?`)
		queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*itemsPerPage)
	}

	if err := r.db.Raw(baseQuery.String(), queryParams...).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].Count

		for i := range data {
			leaveCount, err := r.GetLeaveCountByUser(data[i].DepartmentMemberID)
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
		SELECT COUNT(1)
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
			WHERE DepartmentID = ?`, constant.Inactive, time.Now(), leaveID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *leaveRepository) IsLeaveExistsWithApproval(leaveID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMemberLeaveRequest dmlr
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.ID = ? AND dmlr.IsApproved = 1`, leaveID).
		Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *leaveRepository) GetLeaveCountByUser(departmentMemberID uint) (float64, error) {
	var data []struct {
		ID        uint `gorm:"column:ID"`
		IsFullDay bool `gorm:"column:IsFullDay"`
	}
	var leaveCount float64

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.DepartmentMemberID = ? AND dmlr.IsActive = 1
		AND date(dmlrd.[Date]) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(dmlrd.[Date]) <= date(strftime('%Y-%m', 'now') || '-26')`, departmentMemberID).
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

func (r *leaveRepository) GetLeaveCountByDepartment(departmentID uint) (float64, error) {
	var data []struct {
		ID        uint `gorm:"column:ID"`
		IsFullDay bool `gorm:"column:IsFullDay"`
	}
	var leaveCount float64

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		INNER JOIN DepartmentMember dm ON dm.id = dmlr.DepartmentMemberID AND dm.IsActive = 1
		WHERE dmlr.IsActive = 1 AND dm.DepartmentID = ?
		AND date(dmlrd.[Date]) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(dmlrd.[Date]) <= date(strftime('%Y-%m', 'now') || '-26')`, departmentID).
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

func (r *leaveRepository) GetLeaveCount() (float64, error) {
	var data []struct {
		ID        uint `gorm:"column:ID"`
		IsFullDay bool `gorm:"column:IsFullDay"`
	}
	var leaveCount float64

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.IsActive = 1 
		AND date(dmlrd.[Date]) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(dmlrd.[Date]) <= date(strftime('%Y-%m', 'now') || '-26')`).
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

func (r *leaveRepository) GetApprovedLeaveCount() (float64, error) {
	var data []struct {
		ID        uint `gorm:"column:ID"`
		IsFullDay bool `gorm:"column:IsFullDay"`
	}
	var leaveCount float64

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.IsActive = 1 AND dmlr.IsApproved = 1
		AND date(dmlrd.[Date]) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(dmlrd.[Date]) <= date(strftime('%Y-%m', 'now') || '-26')`).
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

func (r *leaveRepository) GetApprovedLeaveCountByDepartment(departmentID uint) (float64, error) {
	var data []struct {
		ID        uint `gorm:"column:ID"`
		IsFullDay bool `gorm:"column:IsFullDay"`
	}
	var leaveCount float64

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		INNER JOIN DepartmentMember dm ON dm.id = dmlr.DepartmentMemberID AND dm.IsActive = 1
		WHERE dmlr.IsActive = 1 AND dmlr.IsApproved = 1 AND dm.DeparmentMemberID = ?
		AND date(dmlrd.[Date]) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(dmlrd.[Date]) <= date(strftime('%Y-%m', 'now') || '-26')`, departmentID).
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

func (r *leaveRepository) GetApprovedLeaveCountByUser(departmentMemberID uint) (float64, error) {
	var data []struct {
		ID        uint `gorm:"column:ID"`
		IsFullDay bool `gorm:"column:IsFullDay"`
	}
	var leaveCount float64

	if err := r.db.Raw(`
		SELECT dmlrd.ID, dmlrd.IsFullDay 
		FROM DepartmentMemberLeaveRequest dmlr 
		INNER JOIN DepartmentMemberLeaveRequestDate dmlrd 
		ON dmlrd.DepartmentMemberLeaveRequestID = dmlr.ID AND dmlrd.IsActive = 1
		WHERE dmlr.IsActive = 1 AND dmlr.IsApproved = 1 AND dmlr.DeparmentMemberID = ?
		AND date(dmlrd.[Date]) >= date(strftime('%Y-%m', 'now', '-1 month') || '-27')
		AND date(dmlrd.[Date]) <= date(strftime('%Y-%m', 'now') || '-26')`, departmentMemberID).
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
