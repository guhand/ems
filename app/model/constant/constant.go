package constant

type Role int64

const (
	Admin Role = iota + 1
	Manager
	HR
	DepartmentLead
	Employee
)

type Status int

const (
	Inactive Status = iota
	Active
)
