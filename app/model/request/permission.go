package request

import "time"

type RequestPermission struct {
	Reason   string    `json:"reason" binding:"required"`
	Date     time.Time `json:"date" binding:"required"`
	FromTime string    `json:"fromTime" binding:"required"`
	ToTime   string    `json:"toTime" binding:"required"`
}

type UpdatePermissionStatus struct {
	IsApproved bool `json:"isApproved"`
}
