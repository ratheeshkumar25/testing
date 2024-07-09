package middleware

import (
	"errors"
	"fmt"
	"restaurant/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GenerateUserToken generates a JWT token for a user.
func GenerateUsertoken(phone string, userID uint) (string, error) {
	fmt.Println(userID)
	claims := models.UserClaims{
		Phone:  phone,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// AuthenticateUser authenticates a user using the provided JWT token.
func AuthenticateUser(signedStringToken string) (string, uint, error) {
	var userClaims models.UserClaims
	token, err := jwt.ParseWithClaims(signedStringToken, &userClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil // Replace with your secret key
	})

	if err != nil {
		return "", 0, err
	}
	//check the token is valid
	if !token.Valid {
		return "", 0, errors.New("token is not valid")
	}
	//type assert the claims from the token object
	claims, ok := token.Claims.(*models.UserClaims)

	if !ok {
		err = errors.New("couldn't parse claims")
		return "", 0, err
	}
	phone := claims.Phone
	userIDF := float64(claims.UserID)
	fmt.Println(userIDF)
	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token expired")
		return "", 0, err
	}
	userID := uint(userIDF)

	return phone, userID, nil

}

// UserAuthMiddleware is a middleware to authenticate users.
func UserauthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from the request header or other sources
		tokenString := c.GetHeader("Authorization")
		//fmt.Println("Authorization Header",tokenString)

		// Check if token exists
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "User Authorization is missing"})
			return
		}

		// Trim the token to get the actual token string
		//authHeader := strings.TrimSpace(strings.TrimPrefix(tokenString,"Bearer "))
		authHeader := strings.Replace(tokenString, "Bearer ", "", 1)
		phone, userID, err := AuthenticateUser(authHeader)
		if err != nil {
			//fmt.Println("Error authenticating user:", err)
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		c.Set("userID", userID)
		fmt.Println(userID)
		fmt.Println("Authenticated user:", phone)

	}
}
