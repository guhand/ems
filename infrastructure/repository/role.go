package repository

import (
	"ems/app/model/constant"
	"ems/app/model/response"
	"ems/domain"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domain.RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) FetchRoles() ([]response.FetchRoles, error) {
	var data []response.FetchRoles

	if err := r.db.Raw(`
		SELECT ID, [Name]
		FROM [Role]
		WHERE ID <> ?`, constant.Admin).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
