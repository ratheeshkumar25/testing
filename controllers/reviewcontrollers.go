package controllers

import (
	"net/http"
	//"os/user"
	"restaurant/database"
	"restaurant/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ViewReview(c *gin.Context) {
	// Fetch all reviews from the database
	var feedback []models.ReviewModel
	if err := database.DB.Find(&feedback).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	feedBackData := make([]gin.H, len(feedback))
	for i, review := range feedback {
		feedBackData[i] = gin.H{
			"userid":   review.UserID,
			"name":     review.Name,
			"feedback": review.Suggestion,
			"rating":   review.Rating,
		}

	}

	c.JSON(200, gin.H{"message": feedBackData})
}

func Rating(c *gin.Context) {
	// Bind JSON data to feedback model
	var feedback models.ReviewModel
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Get userId from context
	userIDContext, _ := c.Get("userID")
	userID := userIDContext.(uint)

	var booking models.UsersModel
	if err := database.DB.Where("user_id = ?", userID).First(&booking).Error; err == gorm.ErrRecordNotFound {
		c.JSON(200, gin.H{
			"status":  "Success",
			"message": "No booking",
			"data":    nil})
		return
	} else if err != nil {
		c.JSON(404, gin.H{"status": "Failed",
			"message": "Database error",
			"data":    nil})
		return
	}
	feedback.UserID = userID
	//feedback.Name = booking.Username

	if err := database.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	response := gin.H{
		"userID":   feedback.UserID,
		"userName": feedback.Name,
		"review":   feedback.Suggestion,
		"rating":   feedback.Rating,
	}

	c.JSON(http.StatusCreated, response)
}
