package models

import "gorm.io/gorm"

// StaffModel represents a staff member.
type StaffModel struct {
	gorm.Model
	StaffName string `json:"staffname"`
	Role      string `json:"staffrole"`
	Salary    int    `json:"salary"`
	Blocked   bool   `json:"blocked"`
}
