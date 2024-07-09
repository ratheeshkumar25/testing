package controllers

import (
	//"fmt"
	"net/http"
	"restaurant/database"
	"restaurant/models"
	"strconv"

	//"time"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)


// // Create table by admin
func CreateTable(c *gin.Context) {
	var table models.TablesModel
	if err := c.BindJSON(&table); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&table).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{
		"message": "Table created successfully",
		"table": gin.H{
			"tableID":     table.ID,
			"capacity":    table.Capacity,
			"availabilty": table.Availability,
		},
	})
}

// Update the table details
func UpdateTable(c *gin.Context) {
	var table models.TablesModel
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	tableID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid "})
		return
	}
	var existingTable models.TablesModel

	if err := database.DB.First(&existingTable, tableID).Error; err != nil {
		c.JSON(400, gin.H{"error": "Table not found"})
		return
	}
	//update the fileds of existing tablelist
	existingTable.ID = uint(tableID)
	existingTable.Capacity = table.Capacity
	existingTable.Availability = table.Availability

	if err := database.DB.Save(&existingTable).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to update table"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Table updated successfully",
		"table": gin.H{
			"tableID":     table.ID,
			"capacity":    table.Capacity,
			"availabilty": table.Availability,
		},
	})
}

// Remove the table
func RemoveTable(c *gin.Context) {
	tableID := c.Param("id")
	var table models.TablesModel

	if err := database.DB.First(&table, tableID).Error; err != nil {
		c.JSON(404, gin.H{
			"status":  "Failed",
			"message": "tableID not found",
			"data":    err.Error(),
		})
		return
	}

	database.DB.Delete(&table)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Table removed successfully",
		"data": gin.H{
			"tableID":     table.ID,
			"capacity":    table.Capacity,
			"avilability": table.Availability,
		},
	})

}

// Getables retrieves details of all the tables
func GetTables(c *gin.Context) {
	//reterive the tableinformation
	var tables []models.TablesModel
	database.DB.Find(&tables)
	response := make([]gin.H, len(tables))
	for i, tableData := range tables {
		response[i] = gin.H{
			"tableID":  tableData.ID,
			"capacity": tableData.Capacity,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"message": "Table details fetched successfully",
		"data":    response,
	})
}

func GetTable(c *gin.Context) {
	reservationID := c.Param("id")
	var tables models.TablesModel

	if err := database.DB.First(&tables, "id =? ", reservationID).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Failed",
			"message": "Table not found",
			"data":    err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Table found succeesfully",
		"data":    tables,
	})

}
