package request

type CreateDepartment struct {
	Name   string `json:"name" binding:"required"`
	LeadID uint   `json:"leadID" binding:"required"`
}

type UpdateDepartment struct {
	Name   string `json:"name" binding:"required"`
	LeadID uint   `json:"leadID" binding:"required"`
}

type MappUsersToDepartment struct {
	UserIDs []uint `json:"userIDs"`
}

type UnMapEmployee struct {
	UserID uint `json:"userID"`
}
