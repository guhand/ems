package repository

import (
	"ems/app/model/constant"
	"ems/app/model/request"
	"ems/app/model/response"
	"ems/app/model/schema"
	"ems/domain"
	"ems/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateOTP(data *schema.ForgotPasswordOtp) (bool, error) {
	if err := r.db.Exec(`
		INSERT INTO ForgotPasswordOtp
		(CreatedAt, UpdatedAt, IsActive, UserID, Email, Otp, IsUsed)
		VALUES(?, ?, 1, ?, ?, ?, 0)`,
		time.Now(), time.Now(), data.UserID, data.Email, data.Otp).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *userRepository) GetUserByEmail(email string) (*response.FetchUserByEmail, error) {
	var data *response.FetchUserByEmail

	if err := r.db.Raw(`
		SELECT usr.ID, usr.FirstName, usr.LastName, usr.Email, usr.Mobile, usr.token, 
		[Role].ID roleID, [Role].[Name] roleName, usr.CreatedAt, usr.IsActive,
		Usr.[Password], usr.Code, dm.ID AS departmentMemberID, dept.ID AS departmentID, 
		(manager.FirstName || ' ' || manager.LastName) AS manager, manager.ID managerID,
		dept.[Name] AS department, lead.ID AS leadID, (lead.FirstName || ' ' || lead.LastName) AS lead
		FROM [User] usr
		INNER JOIN [Role] ON [Role].ID = usr.RoleID AND [Role].IsActive = 1
		LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.IsActive = 1
		LEFT JOIN Department dept ON dept.ID = dm.departmentID AND dept.IsActive = 1
		LEFT JOIN (
			SELECT leadDM.departmentID, lead.ID, lead.FirstName, lead.LastName
			FROM DepartmentMember leadDM
			INNER JOIN [User] lead ON lead.ID = leadDM.UserID
			WHERE lead.RoleID = ? AND lead.IsActive = 1
		) lead ON lead.departmentID = dept.ID
		LEFT JOIN [User] manager ON manager.id = usr.managerID AND manager.isActive = 1
		WHERE usr.IsActive = 1 AND usr.Email = ?`, constant.DepartmentLead, email).
		Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) GetOTPStatusByUserID(ID uint) (*schema.ForgotPasswordOtp, error) {
	var data *schema.ForgotPasswordOtp

	if err := r.db.Raw(`
		SELECT *
		FROM ForgotPasswordOtp
		WHERE UserID = ?
		ORDER BY CreatedAt DESC LIMIT 1`,
		ID).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) UpdateToken(userID uint, token string) error {
	return r.db.Exec(`
		UPDATE User 
		SET token = ?
		WHERE IsActive = 1 AND ID = ?`, token, userID).Error
}

func (r *userRepository) UpdateOTPStatus(userID uint) error {
	return r.db.Exec(`
		UPDATE ForgotPasswordOtp
		SET isActive = 0, IsUsed = 1
		WHERE UserID = ?`, userID).Error
}

func (r *userRepository) GetUserByID(ID uint) (*response.FetchUserByID, error) {
	var data *response.FetchUserByID

	if err := r.db.Raw(`
		SELECT user.ID userID, user.token token, user.RoleID roleID, 
		dm.ID departmentMemberID, dm.DepartmentID departmentID
		FROM User user
		LEFT JOIN DepartmentMember dm ON dm.UserID = User.ID AND dm.IsActive = 1
		WHERE user.IsActive = 1 AND user.ID = ?`, ID).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) IsEmailExists(email string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User
		WHERE Email = ? AND IsActive = 1`, email).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) IsMobileNumberExists(mobile string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User
		WHERE Mobile = ? AND IsActive = 1`, mobile).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) IsEmailExistsExceptID(id uint, email string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User
		WHERE ID <> ? AND Email = ? AND IsActive = 1`, id, email).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) IsMobileNumberExistsExceptID(id uint, mobile string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User
		WHERE ID <> ? AND Mobile = ? AND IsActive = 1`, id, mobile).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) IsUserExists(id uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User
		WHERE ID = ? AND IsActive = 1`, id).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) IsUserCodeExists(code string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User
		WHERE Code = ? AND IsActive = 1`, code).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) IsUserCodeExistsExceptID(id uint, code string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User
		WHERE ID <> ? AND Code = ? AND IsActive = 1`, id, code).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) CreateUser(req *request.CreateUser, hashedPassword string) error {
	return r.db.Exec(`
		INSERT INTO [User] (
			CreatedAt, UpdatedAt, IsActive, ManagerID, FirstName, LastName, Email, Mobile, 
			Code, RoleID, [Password]
		)
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		time.Now(), time.Now(), constant.Active, 3, // 3 => Manager
		req.FirstName, req.LastName, req.Email, req.Mobile, req.Code,
		req.RoleID, hashedPassword).Error
}

func (r *userRepository) FetchUsers(filters *request.FetchUsers) (*utils.PaginationResponse, error) {
	var (
		data         []response.FetchUsers
		search            = "%" + strings.TrimSpace(filters.Search) + "%"
		itemsPerPage uint = 10
		totalCount   uint = 0
		query        strings.Builder
		queryParams  []interface{}
	)

	query.WriteString(`
		SELECT [User].ID, [User].FirstName, [User].LastName, [User].Email, [User].Code, 
		[User].Mobile, [Role].ID roleID, [Role].[Name] roleName, [User].CreatedAt, 
		Department.ID departmentID, dm.ID AS departmentMemberID,
		[User].UpdatedAt, [User].IsActive, COUNT(*) OVER (PARTITION BY 1) AS userCount, 
		CASE WHEN dm.UserID IS NOT NULL THEN Department.[Name] ELSE 'None' END AS Department, 
		strftime('%Y-%m-%d', DateOfJoining) AS DateOfJoining, Experience, Designation
		FROM [User]
		INNER JOIN [Role] ON [Role].ID = [User].RoleID
		LEFT JOIN UserDetails ud ON [User].ID = ud.UserID AND ud.IsActive = 1
		LEFT JOIN DepartmentMember dm ON dm.UserID = [User].ID AND dm.IsActive = 1
		LEFT JOIN Department ON Department.ID = dm.DepartmentID AND Department.IsActive = 1
		WHERE [User].IsActive = 1`)

	if filters.RoleID > 0 {
		query.WriteString(` AND [User].RoleID = ?`)
		queryParams = append(queryParams, filters.RoleID)
	}

	if filters.DepartmentID > 0 {
		query.WriteString(` AND dm.DepartmentID = ?`)
		queryParams = append(queryParams, filters.DepartmentID)
	}

	if len(filters.Search) > 0 {
		query.WriteString(` AND ([User].Email LIKE ? OR [User].FirstName LIKE ? 
			OR [User].LastName LIKE ?)`)
		queryParams = append(queryParams, search, search, search)
	}

	query.WriteString(` ORDER BY [User].CreatedAt DESC`)

	if filters.Page > 0 {
		query.WriteString(` LIMIT ? OFFSET ?`)
		queryParams = append(queryParams, itemsPerPage, (filters.Page-1)*itemsPerPage)
	}

	if err := r.db.Raw(query.String(), queryParams...).Scan(&data).Error; err != nil {
		return nil, err
	}

	if len(data) > 0 {
		totalCount = data[0].UserCount
	}

	response := *utils.PaginatedResponse(totalCount, filters.Page, data)

	return &response, nil
}

func (r *userRepository) UpdateUser(userID uint, req *request.UpdateUser) error {
	return r.db.Exec(`
				UPDATE [User]
				SET UpdatedAt = ?, FirstName = ?, LastName = ?, Code = ?, Email = ?, Mobile = ?
				WHERE ID = ?`,
		time.Now(), req.FirstName, req.LastName, req.Code,
		req.Email, req.Mobile, userID).
		Error
}

func (r *userRepository) IsUnmappedLeadUser(userID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User usr
		LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.isActive = 1
		WHERE usr.ID = ? AND usr.roleID = ? AND usr.IsActive = 1 AND dm.UserID IS NULL`,
		userID, constant.DepartmentLead).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) IsUnmappedLeadUserIncludeUserID(userID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		WITH tempTable AS(
			SELECT usr.id userID 
			FROM User usr
			LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.IsActive = 1
			WHERE usr.roleID = ? AND usr.IsActive = 1 AND ( dm.UserID = ? OR dm.UserID IS NULL)) 
			
			SELECT COUNT(*)
			FROM tempTable
			WHERE tempTable.userID = ?;`,
		constant.DepartmentLead, userID, userID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) IsUnmappedHRUserIncludeUserID(userID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		WITH tempTable AS(
			SELECT usr.id userID 
			FROM User usr
			LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.IsActive = 1
			WHERE usr.roleID = ? AND usr.IsActive = 1 AND ( dm.UserID = ? OR dm.UserID IS NULL)) 
			
			SELECT COUNT(*)
			FROM tempTable
			WHERE tempTable.userID = ?;`,
		constant.HR, userID, userID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) FetchLastUserCode() (*response.FetchLastUserCode, error) {
	var data *response.FetchLastUserCode

	if err := r.db.Raw(`
		SELECT Code
		FROM [User]
		ORDER BY CreatedAt DESC LIMIT 1`).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) UpdateTokenStatus(userID uint) error {
	return r.db.Exec(`
		UPDATE [User]
		SET Token = NULL
		WHERE ID = ?`, userID).Error
}

func (r *userRepository) UpdatePassword(userId uint, hashedPassword string) error {
	return r.db.Exec(`
			UPDATE User
			SET UpdatedAt = ?, [Password] = ?
			WHERE ID = ?`, time.Now(), hashedPassword, userId).Error
}

func (r *userRepository) FetchUnmappedLeadUsers() ([]response.FetchUnmappedUsers, error) {
	var data []response.FetchUnmappedUsers

	if err := r.db.Raw(`
		SELECT [User].id, ([User].FirstName || ' ' || [User].LastName) AS [name]
		FROM [User]
		LEFT JOIN DepartmentMember dm ON dm.userID = [User].id AND dm.isActive = 1
		WHERE [User].isActive = 1 AND [User].roleID = ? AND dm.userID IS NULL`,
		constant.DepartmentLead).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) FetchUnmappedLeadUserIncludeUserID(req *request.FetchUnmappedLeadUserIncludeUserID) ([]response.FetchUnmappedUsers, error) {
	var data []response.FetchUnmappedUsers

	if err := r.db.Raw(`
			SELECT DISTINCT usr.id, (usr.FirstName || ' ' || usr.LastName) AS [name]
			FROM User usr
			LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.IsActive = 1
			WHERE usr.roleID = ? AND usr.IsActive = 1 AND ( dm.UserID = ? OR dm.UserID IS NULL)`,
		constant.DepartmentLead, req.UserID).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) GetUnmappedEmployeesCount(userIDs []uint) (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User usr
		LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.isActive = 1
		WHERE usr.ID IN ? AND usr.roleID = ? AND usr.IsActive = 1 AND dm.UserID IS NULL`,
		userIDs, constant.Employee).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) FetchUnmappedEmployeeUsers() ([]response.FetchUnmappedUsers, error) {
	var data []response.FetchUnmappedUsers

	if err := r.db.Raw(`
		SELECT usr.id, (usr.FirstName || ' ' || usr.LastName) AS [name]
		FROM User usr
		LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.isActive = 1
		WHERE usr.roleID = ? AND usr.IsActive = 1 AND dm.UserID IS NULL`,
		constant.Employee).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) IsDepartmentUserExists(userID uint) (*response.FetchDepartmentUserCountAndRoleID, error) {
	var data *response.FetchDepartmentUserCountAndRoleID

	if err := r.db.Raw(`
		SELECT COUNT(*) AS [count], usr.RoleID roleID, dm.DepartmentID departmentID
		FROM User usr
		INNER JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.isActive = 1
		WHERE usr.ID = ? AND usr.IsActive = 1`, userID).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) RemoveUser(userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			UPDATE User
			SET IsActive = ?, DeletedAt = ?
			WHERE ID = ?`, constant.Inactive, time.Now(), userID).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE DepartmentMember
			SET IsActive = ?, DeletedAt = ?
			WHERE UserID = ?`, constant.Inactive, time.Now(), userID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *userRepository) IsMappedLeadUser(userID uint) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User usr
		LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.isActive = 1
		WHERE usr.ID = ? AND usr.roleID IN ? AND usr.IsActive = 1 AND dm.UserID IS NOT NULL`,
		userID, []interface{}{constant.DepartmentLead, constant.HR}).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) UpdateUserDetails(req *request.UpdateUserDetails) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		layout := "2006-01-02"
		dob := req.DOB.Format(layout)
		doj := req.DateOfJoining.Format(layout)

		var count int64
		if err := tx.Raw(`
				SELECT COUNT(*) 
				FROM UserDetails
				WHERE IsActive = 1 AND UserID = ?`, req.UserID).Scan(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			if err := tx.Exec(`
				UPDATE UserDetails
				SET UpdatedAt = ?, DateOfJoining = strftime('%Y-%m-%d', ?), Experience = ?, Designation = ?, 
				DOB = strftime('%Y-%m-%d', ?), PanNumber = ?, AadharNumber = ?, 
				BankAccountNumber = ?, IfscCode = ?, City = ?, [Address] = ?, 
				Degree = ?, College = ?
				WHERE UserID = ?`,
				time.Now(), doj, req.Experience, req.Designation, dob, req.PanNumber,
				req.AadharNumber, req.BankAccountNumber, req.IfscCode,
				req.City, req.Address, req.Degree, req.College, req.UserID).
				Error; err != nil {
				return err
			}
		} else {
			if err := tx.Exec(`
				INSERT INTO UserDetails
				(CreatedAt, UpdatedAt, UserID, IsActive, DateOfJoining, DOB, Experience, Designation,
				 PanNumber, AadharNumber, BankAccountNumber, IfscCode, City, [Address], 
				Degree, College)
				VALUES(?, ?, ?, ?, strftime('%Y-%m-%d', ?), strftime('%Y-%m-%d', ?), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				time.Now(), time.Now(), req.UserID, constant.Active, doj, dob, req.Experience,
				req.Designation, req.PanNumber, req.AadharNumber, req.BankAccountNumber,
				req.IfscCode, req.City, req.Address, req.Degree, req.College).
				Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *userRepository) IsAadharNumberExistsExceptID(id uint, aadharNumber string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM UserDetails
		WHERE UserID <> ? AND AadharNumber = ?`, id, aadharNumber).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) IsPanNumberExistsExceptID(id uint, panNumber string) (bool, error) {
	var count int64

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM UserDetails
		WHERE UserID <> ? AND PanNumber = ?`, id, panNumber).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) FetchUserDetails(req *request.FetchUserDetails) (*response.FetchUserDetails, error) {
	var data *response.FetchUserDetails

	if err := r.db.Raw(`
		SELECT usr.ID, usr.FirstName, usr.LastName, usr.Email, usr.Mobile, 
		strftime('%Y-%m-%d', DOB) AS DOB, strftime('%Y-%m-%d', DateOfJoining) AS DateOfJoining,
		PanNumber, AadharNumber, Experience, Designation,
		BankAccountNumber, IfscCode, City, [Address], Degree, College
		FROM [User] usr
		LEFT JOIN UserDetails ud ON usr.ID = ud.UserID AND ud.IsActive = 1
		WHERE usr.ID = ?`,
		req.UserID).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil

}

func (r *userRepository) UploadFiles(userID uint, filepaths []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, filepath := range filepaths {
			if err := tx.Exec(`
				INSERT INTO UserDocument
				(CreatedAt, UpdatedAt, UserID, FilePath)
				VALUES(?, ?, ?, ?)`, time.Now(), time.Now(), userID, filepath).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepository) FetchUnmappedHRUserIncludeUserID(req *request.FetchUnmappedLeadUserIncludeUserID) ([]response.FetchUnmappedUsers, error) {
	var data []response.FetchUnmappedUsers

	if err := r.db.Raw(`
			SELECT DISTINCT usr.id, (usr.FirstName || ' ' || usr.LastName) AS [name]
			FROM User usr
			LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.IsActive = 1
			WHERE usr.roleID = ? AND usr.IsActive = 1 AND (dm.UserID = ? OR dm.UserID IS NULL)`,
		constant.HR, req.UserID).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) FetchUnmappedHRUsers() ([]response.FetchUnmappedUsers, error) {
	var data []response.FetchUnmappedUsers

	if err := r.db.Raw(`
		SELECT usr.id, (usr.FirstName || ' ' || usr.LastName) AS [name]
		FROM User usr
		LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.IsActive = 1
		WHERE usr.RoleID = ? AND usr.IsActive = 1 AND dm.UserID IS NULL`,
		constant.HR).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) GetUnmappedHRsCount(userIDs []uint) (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User usr
		LEFT JOIN DepartmentMember dm ON dm.UserID = usr.ID AND dm.isActive = 1
		WHERE usr.ID IN ? AND usr.RoleID = ? AND usr.IsActive = 1 AND dm.UserID IS NULL`,
		userIDs, constant.HR).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) FetchFilePathsByUserID(userID uint) ([]response.FetchUploadedDocumentPaths, error) {
	var data []response.FetchUploadedDocumentPaths

	if err := r.db.Raw(`
		SELECT ud.UserID, GROUP_CONCAT(filePath) AS filePaths
		FROM UserDocument ud
		WHERE ud.UserID = ? AND ud.IsActive = 1
		GROUP BY ud.userID`, userID).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *userRepository) RemoveFile(filePath string) error {
	return r.db.Exec(`
		UPDATE UserDocument
		SET IsActive = ?, DeletedAt = ?
		WHERE FilePath = ?`, constant.Inactive, time.Now(), filePath).Error
}

func (r *userRepository) GetUserCount() (int, error) {
	var count int

	if err := r.db.Raw(`
		SELECT COUNT(*)
		FROM User usr
		WHERE usr.IsActive = 1 AND usr.RoleID <> ?`,
		constant.Admin).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
