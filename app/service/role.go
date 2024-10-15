package service

import (
	"ems/app/model/response"
	"ems/domain"
)

type roleService struct {
	roleRepository domain.RoleRepository
}

func NewRoleService(roleRepository domain.RoleRepository) domain.RoleService {
	return &roleService{roleRepository}
}

func (s *roleService) FetchRoles() ([]response.FetchRoles, error) {
	data, err := s.roleRepository.FetchRoles()

	if err != nil {
		return nil, err
	}

	return data, nil
}
