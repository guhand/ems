package schema

import (
	"time"

	"gorm.io/gorm"
)

type BaseGorm struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
	IsActive  bool `gorm:"default:true"`
}

type Role struct {
	BaseGorm
	Name  string `gorm:"not null"`
	Users []User
}
type User struct {
	BaseGorm
	FirstName           string `gorm:"not null"`
	LastName            string `gorm:"not null"`
	Email               string `gorm:"not null"`
	Mobile              string `gorm:"not null"`
	Code                string `gorm:"not null"`
	Password            string `gorm:"not null"`
	Token               *string
	RoleID              uint `gorm:"not null"`
	Role                Role
	ManagerID           *uint `gorm:"foreignKey:ManagerID"`
	Manager             *User `gorm:"foreignKey:ManagerID"`
	UserDetails         []UserDetails
	UserDocuments       []UserDocument
	ForgotPasswordOtps  []ForgotPasswordOtp
	DepartmentMembers   []DepartmentMember
	ApprovedLeaves      []DepartmentMemberLeaveRequest      `gorm:"foreignKey:ApprovedBy"`
	ApprovedPermissions []DepartmentMemberPermissionRequest `gorm:"foreignKey:ApprovedBy"`
	ApprovedNotices     []UserNotice                        `gorm:"foreignKey:ApprovedBy"`
}

type UserDetails struct {
	BaseGorm
	UserID            uint `gorm:"not null"`
	User              User
	DateOfJoining     time.Time `gorm:"not null;type:date"`
	Designation       string    `gorm:"not null"`
	Experience        uint      `gorm:"not null"`
	DOB               time.Time `gorm:"not null;type:date"`
	AadharNumber      string    `gorm:"not null"`
	PanNumber         string    `gorm:"not null"`
	BankAccountNumber string    `gorm:"not null"`
	IfscCode          string    `gorm:"not null"`
	Address           string    `gorm:"not null"`
	City              string    `gorm:"not null"`
	Degree            string    `gorm:"not null"`
	College           string    `gorm:"not null"`
}

type UserDocument struct {
	BaseGorm
	UserID   uint `gorm:"not null"`
	User     User
	FilePath string `gorm:"not null"`
}

type ForgotPasswordOtp struct {
	BaseGorm
	UserID uint   `json:"userID" gorm:"not null"`
	User   User   `json:"user"`
	Email  string `json:"email" gorm:"not null"`
	Otp    string `json:"otp" gorm:"not null"`
	IsUsed bool   `json:"isUsed" gorm:"default:false"`
}

type Department struct {
	BaseGorm
	Name              string `gorm:"not null"`
	DepartmentMembers []DepartmentMember
}

type DepartmentMember struct {
	BaseGorm
	DepartmentID                       uint `gorm:"not null"`
	Department                         Department
	UserID                             uint `gorm:"not null"`
	User                               User
	DepartmentMemberLeaveRequests      []DepartmentMemberLeaveRequest
	DepartmentMemberPermissionRequests []DepartmentMemberPermissionRequest
	DepartmentMemberNotices            []UserNotice
}

type DepartmentMemberLeaveRequest struct {
	BaseGorm
	DepartmentMemberID                uint `gorm:"not null"`
	DepartmentMember                  DepartmentMember
	ShiftType                         uint   `gorm:"default:1"`
	Reason                            string `gorm:"not null"`
	IsApproved                        *bool
	ApprovedAt                        *time.Time
	ApprovedBy                        *uint
	ApprovedUser                      *User `gorm:"foreignKey:ApprovedBy"`
	DepartmentMemberLeaveRequestDates []DepartmentMemberLeaveRequestDate
}

type DepartmentMemberLeaveRequestDate struct {
	BaseGorm
	DepartmentMemberLeaveRequestID uint `gorm:"not null"`
	DepartmentMemberLeaveRequest   DepartmentMemberLeaveRequest
	Date                           time.Time `gorm:"not null;type:date"`
	IsFullDay                      bool
	SessionType                    *uint `gorm:"default:0"`
}

type DepartmentMemberPermissionRequest struct {
	BaseGorm
	DepartmentMemberID uint `gorm:"not null"`
	DepartmentMember   DepartmentMember
	ShiftType          uint      `gorm:"default:1"`
	Reason             string    `gorm:"not null"`
	Date               time.Time `gorm:"not null;type:date"`
	FromTime           string    `gorm:"not null"`
	ToTime             string    `gorm:"not null"`
	IsApproved         *bool
	ApprovedAt         *time.Time
	ApprovedBy         *uint
	ApprovedUser       *User `gorm:"foreignKey:ApprovedBy"`
}

type UserNotice struct {
	BaseGorm
	DepartmentMemberID uint `gorm:"not null"`
	DepartmentMember   DepartmentMember
	Remarks            string `gorm:"not null"`
	NoticeEndDate      *time.Time
	IsApproved         bool  `gorm:"default:false"`
	ApprovedBy         *uint `json:"approvedBy"`
	ApprovedUser       *User `gorm:"foreignKey:ApprovedBy"`
}
