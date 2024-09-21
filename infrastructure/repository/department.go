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

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &departmentRepository{db}
}

func (r *departmentRepository) IsDepartmentExists(id uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM Department
		WHERE ID = ? AND IsActive = 1`, id).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *departmentRepository) CreateDepartment(req *request.CreateDepartment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			INSERT INTO Department
			(CreatedAt, UpdatedAt, IsActive, [Name])
			VALUES(?, ?, 1, ?)`, time.Now(), time.Now(), req.Name).Error; err != nil {
			return err
		}

		var departmentID uint

		if err := tx.Raw(`
			SELECT ID
			FROM Department
			ORDER BY CreatedAt DESC LIMIT 1`).Scan(&departmentID).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			INSERT INTO DepartmentMember (CreatedAt, UpdatedAt, DepartmentID, UserID)
			VALUES(?, ?, ?, ?)`, time.Now(), time.Now(), departmentID, req.LeadID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *departmentRepository) IsDepartmentNameExists(name string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM Department
		WHERE [Name] = ?`, name).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *departmentRepository) FectchDepartments(filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	var data []response.FectchDepartments
	var search = "%" + strings.TrimSpace(filters.Search) + "%"
	var itemsPerPage uint = 10
	var totalCount uint = 0
	var baseQuery strings.Builder
	var queryParams []interface{}

	baseQuery.WriteString(`
		SELECT dept.ID, dept.Name, usr.ID leadID, (usr.FirstName || ' ' || usr.LastName) AS leadName, 
		COUNT(1) OVER (PARTITION BY 1) AS departmentCount, dept.CreatedAt, 
		dept.UpdatedAt, dept.IsActive  
		FROM Department dept
		INNER JOIN DepartmentMember dm ON dm.DepartmentID = dept.ID AND dm.isActive = 1
		INNER JOIN [User] usr ON usr.ID = dm.userID AND usr.isActive = 1 AND usr.roleID IN ?
		WHERE dept.IsActive = 1`)

	queryParams = append(queryParams, []interface{}{constant.DepartmentLead, constant.HR})

	if len(filters.Search) > 0 {
		baseQuery.WriteString(` AND dept.Name LIKE ?`)
		queryParams = append(queryParams, search)
	}

	baseQuery.WriteString(` ORDER BY dept.CreatedAt DESC`)

	if filters.Page > 0 {
		baseQuery.WriteString(` LIMIT ? OFFSET ?`)
		queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*itemsPerPage)
	}

	if err := r.db.Raw(baseQuery.String(), queryParams...).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].DepartmentCount
	}

	response := *utils.PaginatedResponse(totalCount, filters.Page, data)

	return &response, nil
}

func (r *departmentRepository) IsDepartmentNameExistsExceptID(id uint, name string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM Department
		WHERE ID <> ? AND Name = ?`, id, name).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *departmentRepository) UpdateDepartment(id uint, req *request.UpdateDepartment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			UPDATE Department
			SET UpdatedAt = ?, [Name] = ?
			WHERE ID = ?`,
			time.Now(), req.Name, id).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE DepartmentMember
			SET isActive = 0, DeletedAt = ?
			WHERE departmentID = ? AND userID <> ?`,
			time.Now(), id, req.LeadID).Error; err != nil {
			return err
		}

		// Insert into DepartmentMember if not exists
		if err := tx.Exec(`
			INSERT INTO DepartmentMember (CreatedAt, UpdatedAt, DepartmentID, UserID)
			SELECT ?, ?, ?, ?
			WHERE NOT EXISTS (
				SELECT 1 FROM DepartmentMember WHERE isActive = 1 AND departmentID = ? AND userID = ?
			)`,
			time.Now(), time.Now(), id, req.LeadID, id, req.LeadID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *departmentRepository) MappUsersToDepartment(departmentID uint, req *request.MappUsersToDepartment) error {
	var placeholders []string
	var args []interface{}
	now := time.Now()

	query := "INSERT INTO DepartmentMember (CreatedAt, UpdatedAt, DepartmentID, UserID) VALUES "

	for _, userID := range req.UserIDs {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		args = append(args, now, now, departmentID, userID)
	}

	query += strings.Join(placeholders, ", ")

	if err := r.db.Exec(query, args...).Error; err != nil {
		return err
	}

	return nil
}

func (r *departmentRepository) FetchDepartmentMembers(departmentID uint, filters *request.CommonRequest) (*utils.PaginationResponse, error) {
	var (
		itemsPerPage uint = 10
		totalCount   uint = 0
		search            = "%" + strings.TrimSpace(filters.Search) + "%"
		baseQuery    strings.Builder
		queryParams  []interface{}
		data         []response.FetchDepartmentMembers
	)

	baseQuery.WriteString(`
		SELECT usr.ID userID, (usr.FirstName || ' ' || usr.LastName) AS userName, 
		COUNT(1) OVER (PARTITION BY 1) AS departmentMemberCount, usr.CreatedAt, [Role].[name] AS [role]
		FROM Department dept
		INNER JOIN DepartmentMember dm ON dm.DepartmentID = dept.ID AND dm.IsActive = 1
		INNER JOIN [User] usr ON usr.ID = dm.userID AND usr.IsActive = 1
		INNER JOIN [Role] ON usr.roleID = [Role].id AND [Role].IsActive = 1
		WHERE dept.IsActive = 1 AND dept.ID = ?`)

	queryParams = append(queryParams, departmentID)

	if len(filters.Search) > 0 {
		baseQuery.WriteString(` AND (usr.FirstName LIKE ? OR usr.LastName LIKE ?)`)
		queryParams = append(queryParams, search, search)
	}

	baseQuery.WriteString(` ORDER BY dept.CreatedAt DESC`)

	if filters.Page > 0 {
		baseQuery.WriteString(` LIMIT ? OFFSET ?`)
		queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*itemsPerPage)
	}

	if err := r.db.Raw(baseQuery.String(), queryParams...).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].DepartmentMemberCount
	}

	response := *utils.PaginatedResponse(totalCount, filters.Page, data)

	return &response, nil
}

func (r *departmentRepository) UnMapEmployee(req *request.UnMapEmployee) error {
	return r.db.Exec(`
		UPDATE DepartmentMember
		SET IsActive = 0, DeletedAt = ?
		WHERE UserID = ?`, time.Now(), req.UserID).Error
}

func (r *departmentRepository) IsDepartmentMemberExists(id uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMember
		WHERE ID = ? AND IsActive = 1`, id).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *departmentRepository) FetchDepartmentLead(departmentID uint) (*response.FetchDepartmenLead, error) {
	var data response.FetchDepartmenLead

	if err := r.db.Raw(`
		SELECT usr.id leadID, (usr.FirstName || ' ' || usr.LastName) AS leadName
		FROM User usr
		INNER JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.isActive = 1
		INNER JOIN Department dept ON dept.ID = dm.departmentID AND dept.isActive = 1
		WHERE usr.roleID = ? AND usr.IsActive = 1 AND dm.departmentID = ?`,
		constant.DepartmentLead, departmentID).Scan(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *departmentRepository) RemoveDepartment(departmentID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			UPDATE Department
			SET IsActive = ?, DeletedAt = ?
			WHERE ID = ?`, constant.Inactive, time.Now(), departmentID).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE DepartmentMember
			SET IsActive = ?, DeletedAt = ?
			WHERE DepartmentID = ?`, constant.Inactive, time.Now(), departmentID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *departmentRepository) GetDepartmentCount() (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM Department
		WHERE IsActive = 1`).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *departmentRepository) GetDepartmentMemberCount(departmentID uint) (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(1)
		FROM DepartmentMember
		WHERE IsActive = 1 AND DepartmentID = ?`, departmentID).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
