package models

import "github.com/dgrijalva/jwt-go"


// AdminModel represents the admin entity.
type AdminModel struct {
	AdminID  int    `gorm:"primaryKey" `
	Username string `json:"username"`
	Password string `json:"password"`
}
// AdminClaims represents the claims of the admin JWT token.
type AdminClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
