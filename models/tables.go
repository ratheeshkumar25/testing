package models

import (
	"time"

	"gorm.io/gorm"
)

// TablesModel represents a table entity.
type TablesModel struct {
	gorm.Model
	Capacity     int  `json:"capacity" validate:"required"`
	Availability bool `json:"availability" validate:"required" `
}

// ReservationModels represents a reservation entity.
type ReservationModels struct {
	gorm.Model
	Date          time.Time `json:"date" gorm:"column:date"`
	TableID       uint      `json:"tableID"`
	Email         string    `json:"email"`
	NumberOfGuest int       `json:"numberofGuest"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	UserID        uint      `json:"userID"`
	StaffID       uint      `gorm:"not null"`
}

// BeforeCreate hook to set the Date field to the current date before creating a new record.
func (reservation *ReservationModels) BeforeCreate(tx *gorm.DB) (err error) {
	reservation.Date = time.Now()
	return nil
}
