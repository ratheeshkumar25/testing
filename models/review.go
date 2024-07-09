package models

import (


	"gorm.io/gorm"
)

// ReviewModel represents a review entity.
type ReviewModel struct {
	gorm.Model
	UserID     uint 	 `json:"user_id"`
	Name       string    `json:"name"`
	Suggestion string    `json:"suggestion"`
	Rating     int       `json:"rating"`
}
