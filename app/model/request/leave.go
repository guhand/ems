package request

import "time"

type RequestLeave struct {
	Reason string `json:"reason" binding:"required"`
	Dates  []Date `json:"dates"`
}

type Date struct {
	Date        time.Time `json:"date" binding:"required"`
	IsFullDay   bool      `json:"isFullDay"`
	SessionType uint      `json:"sessionType"`
}

type UpdateLeaveStatus struct {
	IsApproved bool `json:"isApproved"`
}
