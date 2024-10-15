package domain

import "ems/app/model/response"

type RoleService interface {
	FetchRoles() ([]response.FetchRoles, error)
}

type RoleRepository interface {
	FetchRoles() ([]response.FetchRoles, error)
}
