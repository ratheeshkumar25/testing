package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"restaurant/models"
	"strings"
	"time"
)

var (
	jwtKey = []byte("sdjgertoweipskfrtqw")
)

// GenerateAdminToken generates a JWT token for an admin.
func GenerateAdminToken(username string) (string, error) {
	//set the exipartion time for token
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the claims for the token
	claims := &models.AdminClaims{
		Username:       username,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// AdminAuthentication verifies the JWT token and returns the admin's username.
func AdminAuthentication(tokenString string) (string, error) {
	//fmt.Println("Received token",tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &models.AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		//fmt.Println("Error parsing/validating token:", err)
		return "", err
	}

	if claims, ok := token.Claims.(*models.AdminClaims); ok && token.Valid {
		//fmt.Println("Token claims:", claims)
		return claims.Username, nil
	}

	return "", errors.New("invalid token")
}

// AdminAuthMiddleware is a middleware function for admin authentication.
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		// fmt.Println("Authorization Header:", tokenString)
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing the authorization header"})
			return
		}

		authHeader := strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer"))
		// fmt.Println(authHeader)

		username, err := AdminAuthentication(authHeader)
		// fmt.Println(username)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
