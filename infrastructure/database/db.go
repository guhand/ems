package database

import (
	"ems/app/model"
	"ems/app/model/constant"
	"ems/app/model/schema"
	"ems/infrastructure/config"
	"ems/utils"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gorm_schema "gorm.io/gorm/schema"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.Config.DbDsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}), NamingStrategy: gorm_schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := migrateSchema(db); err != nil {
		return nil, err
	}
	if err := initData(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(&schema.Role{}, &schema.User{}, &schema.ForgotPasswordOtp{},
		&schema.Department{}, &schema.DepartmentMember{}, &schema.UserNotice{},
		&schema.UserDocument{}, &schema.DepartmentMemberLeaveRequest{}, &schema.UserDetails{},
		&schema.DepartmentMemberLeaveRequestDate{}, &schema.DepartmentMemberPermissionRequest{})
}

func initData(db *gorm.DB) error {
	if err := initRole(db); err != nil {
		return err
	}

	if err := initAdmin(db); err != nil {
		return err
	}

	if err := initHRDepartment(db); err != nil {
		return err
	}

	return nil
}

func initRole(db *gorm.DB) error {
	var count int64

	if err := db.Raw(`SELECT COUNT(1) FROM [Role]`).Scan(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		for _, role := range model.Roles {
			if err := db.Exec(`
				INSERT INTO [Role]
				(CreatedAt, UpdatedAt, IsActive, [Name])
				VALUES(?, ?, 1, ?)`, time.Now(), time.Now(), role).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func initAdmin(db *gorm.DB) error {
	var count int64

	if err := db.Raw(`
		SELECT COUNT(1) 
		FROM User`).Scan(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		now := time.Now()

		hashedPassword, err := utils.HashPassword("123")

		if err != nil {
			return err
		}

		if err := db.Exec(`
			INSERT INTO [User]
			(CreatedAt, UpdatedAt, IsActive, FirstName, LastName, Code, Email, 
			Mobile, RoleID, [Password])
			VALUES
			(?, ?, ?, "Tamil", "M", "A001","tamilselvammca2021@gmail.com", "8438324027", ?, ?),
			(?, ?, ?, "Guhan", "D", "A000","guhand@sternadevices.in", "8438379027", ?, ?),
			(?, ?, ?, "Sterna", "Manager", "M000","ma@test.hf", "9839749324", ?, ?),
			(?, ?, ?, "Sterna", "HR", "H001","ma@t4est.hf", "9839744324", ?, ?)`,
			now, now, constant.Active, constant.Admin,
			hashedPassword, now, now, constant.Active,
			constant.Admin, hashedPassword, now, now, constant.Active,
			constant.Manager, hashedPassword, now, now,
			constant.Active, constant.HR, hashedPassword).Error; err != nil {
			return err
		}

	}
	return nil
}

func initHRDepartment(db *gorm.DB) error {
	var count int64

	if err := db.Raw(`
		SELECT COUNT(1) 
		FROM Department`).Scan(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(`
			INSERT INTO Department
			(CreatedAt, UpdatedAt, IsActive, [Name])
			VALUES(?, ?, 1, ?)`, time.Now(), time.Now(), "Human Resource").Error; err != nil {
				return err
			}

			var departmentID uint

			if err := tx.Raw(`
			SELECT ID
			FROM Department
			ORDER BY CreatedAt DESC LIMIT 1`).Scan(&departmentID).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
			INSERT INTO DepartmentMember (CreatedAt, UpdatedAt, DepartmentID, UserID)
			VALUES(?, ?, ?, ?)`, time.Now(), time.Now(), departmentID, 4).Error; err != nil { // 4 => HR
				return err
			}
			return nil
		})
	}

	return nil
}
