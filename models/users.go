package models

import "github.com/dgrijalva/jwt-go"

// UsersModel represents a user entity.
type UsersModel struct {
	UserID   int    `gorm:"primaryKey;autoIncrement"`
	Phone    string `json:"phone" validate:"required"`
	// Username string `json:"username"`
}

// VerifyOTP represents the structure for verifying OTP.
type VerifyOTP struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Otp      string `json:"otp"`
}

// UserClaims represents the claims of the user JWT token.
type UserClaims struct {
	UserID uint
	Phone  string `json:"phone"`
	jwt.StandardClaims
}
