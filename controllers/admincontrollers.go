package controllers

import (
	"fmt"
	"net/http"
	"restaurant/database"
	"restaurant/middleware"
	"restaurant/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Admin Login
func AdminLogin(c *gin.Context) {
	var admin models.AdminModel
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Retrieve admin from the database
	var dbAdmin models.AdminModel
	if err := database.DB.Where("username = ?", admin.Username).First(&dbAdmin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check if the stored password is hashed
	if len(dbAdmin.Password) > 0 && dbAdmin.Password[0] == '$' {
		// Password is hashed, compare with bcrypt
		if err := bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
	} else {
		// Password is plaintext, compare directly
		if dbAdmin.Password != admin.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		// Hash the plaintext password and update it in the database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		dbAdmin.Password = string(hashedPassword)
		if err := database.DB.Save(&dbAdmin).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}
	}
	fmt.Println("Here is dbadmin", dbAdmin)
	// Generate JWT token
	token, err := middleware.GenerateAdminToken(admin.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Authentication successful, return JWT token
	c.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})

}

// Admin Logout
func AdminLogout(c *gin.Context) {
	//Respond with Successful Logout

	
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Logout"})
}
