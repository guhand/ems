package request

type CreateDepartment struct {
	Name   string `json:"name" binding:"required"`
	LeadID uint   `json:"leadID" binding:"required"`
}

type UpdateDepartment struct {
	CreateDepartment
}

type MappUsersToDepartment struct {
	UserIDs []uint `json:"userIDs"`
}

type UnMapUser struct {
	UserID uint  `json:"userID" binding:"required"`
	LeadID *uint `json:"leadID"`
}
