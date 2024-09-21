package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

type Scheduler struct {
	DB *gorm.DB
}

func (s *Scheduler) InitScheduler() {
	// Initialize the gocron scheduler
	scheduler := gocron.NewScheduler(time.Local)

	// Schedule the job to run every day at midnight
	_, err := scheduler.Every(1).Day().At("00:00").Do(s.removeUsers)
	if err != nil {
		log.Fatalf("Failed to schedule job: %v", err)
	}

	// Start the scheduler asynchronously
	scheduler.StartAsync()
}

func (s *Scheduler) removeUsers() {
	now := time.Now()

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Fetch users whose notice period has ended
		var departmentMemberIDs []uint
		if err := tx.Raw(`
			SELECT DepartmentMemberID 
			FROM UserNotice 
			WHERE NoticeEndDate < ?`, now).Scan(&departmentMemberIDs).Error; err != nil {
			return err
		}

		if len(departmentMemberIDs) > 0 {
			// Mark the corresponding department members as inactive and set DeletedAt
			if err := tx.Exec(`
				UPDATE DepartmentMember 
				SET IsActive = 0, DeletedAt = ? 
				WHERE ID IN ?`, now, departmentMemberIDs).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				UPDATE [User] 
				SET IsActive = 0, DeletedAt = ? 
				WHERE ID IN (
					SELECT UserID FROM DepartmentMember WHERE ID IN ?
				)`, now, departmentMemberIDs).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	} else {
		fmt.Println("Users removed successfully")
	}
}
