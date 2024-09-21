package request

type CreateNotice struct {
	Remarks string `json:"remarks" binding:"required"`
}

type ApproveNotice struct {
	ServeDays int `json:"serveDays" binding:"required"`
}
