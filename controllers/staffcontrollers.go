package controllers

import (
	"restaurant/database"
	"restaurant/models"
	"github.com/gin-gonic/gin"
)

// Getstaff details
func GetStaff(c *gin.Context) {
	var staff []models.StaffModel
	if err := database.DB.Find(&staff).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//Prepare Staff data for response,including pnly desire fields
	staffData := make([]gin.H, len(staff))
	for i, staffMember := range staff {
		staffData[i] = gin.H{
			"staffID":     staffMember.ID,
			"staffName":   staffMember.StaffName,
			"designation": staffMember.Role,
		}
	}
	c.JSON(200, gin.H{"staff": staffData})

}

// GetStaffByID retrieve a staff member by ID
func GetStaffByIDs(c *gin.Context) {
	var staff models.StaffModel
	staffID := c.Param("id")

	if err := database.DB.First(&staff, staffID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Staff not found"})
		return
	}
	//Prepare response for staff data , including only desire fields
	respone := gin.H{
		"message": "staff details fetched successfully",
		"data": gin.H{
			"staffID":     staff.ID,
			"staffName":   staff.StaffName,
			"Designation": staff.Role,
		},
	}

	c.JSON(200, respone)

}

// Add new staff member to restaurant
func AddStaff(c *gin.Context) {
	var staff models.StaffModel
	if err := c.BindJSON(&staff); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&staff).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"Status":    "Success",
		"Message":   "Successfully added staff details",
		"staffID":   staff.ID,
		"staffname": staff.StaffName,
		"Role":      staff.Role,
	})
}

// Update a staff member
func UpdateStaff(c *gin.Context) {
	var staff models.StaffModel
	staffID := c.Param("id")

	if err := database.DB.First(&staff, staffID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Staff not found"})
		return
	}
	if err := c.BindJSON(&staff); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Staff details updated successfully",
		"data":    staff,
	})

}

// Remove a staff member
func RemoveStaff(c *gin.Context) {
	staffID := c.Param("id")
	var staff models.StaffModel

	if err := database.DB.First(&staff, staffID).Error; err != nil {
		c.JSON(404, gin.H{
			"status":  "Failed",
			"message": "Staff id not found",
			"data":    err.Error(),
		})
		return
	}
	database.DB.Delete(&staff)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "staff details removed successfully",
	})
}

// func StaffAssignTable(c *gin.Context){
// 	var request struct{
// 		StaffID uint `json:"staffID"`
// 		TableID uint `json:"tableID"`
// 	}
// 	if err:= c.ShouldBindJSON(&request); err != nil{
// 		c.JSON(400,gin.H{"error":err.Error()})
// 		return
// 	}

// 	//Fetch the staff details
// 	var staff models.StaffModel
// 	if err := database.DB.First(&staff,request.StaffID).Error; err != nil{
// 		c.JSON(404,gin.H{"error":"Staff not found"})
// 		return
// 	}
// 	//staff.TableID = request.TableID
// 	if err := database.DB.Save(&staff).Error; err!=nil{
// 		c.JSON(500,gin.H{"error":"Failed to assign table to staff"})
// 	}
// 	c.JSON(200,gin.H{
// 		"status":"Success",
// 		"message":"Table assigned to staff successfully",
// 		"data":gin.H{
// 			"staff":staff,
// 			"tableID":request.TableID,
// 		},
// 	})
// }

// func fetchStaffIDByTableID(tableID int) (int, error) {
// 	var staff models.StaffModel
// 	if err := database.DB.Where("table_id = ?", tableID).First(&staff).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return 0, fmt.Errorf("no staff found for table with ID %d", tableID)
// 		}
// 		return 0, fmt.Errorf("failed to fetch staff: %v", err)
// 	}

// 	return int(staff.ID), nil
// }
