package controllers

import (
	"fmt"
	"net/http"
	"restaurant/database"
	"restaurant/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// User can access whole menu list
func GetMenuList(c *gin.Context) {
	var menus []models.MenuModel
	database.DB.Find(&menus)
	//fmt.Println(menus)

	//Prpare menu data for response ,including only desire fields
	menuData := make([]gin.H,len(menus))
	for i , menuItem := range menus{
		menuData[i] = gin.H{
			"menuID":menuItem.ID,
			"name":menuItem.Name,
			"category":menuItem.Category,
			"price":menuItem.Price,
			"duration":menuItem.Duration,
		}
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Menu details fetched successfully",
		"menulist": menuData,
	})
}

// Access to particular menu for user to check with food id
func GetMenu(c *gin.Context) {
	//Reterive the food_id parameter from the URL
	foodID := c.Param("id")

	//Declare a slice to hold menuitems
	var menuItem models.MenuModel

	//Find all items matched with that match the food_id
	if err := database.DB.First(&menuItem, "ID = ?", foodID).Error; err != nil {
		c.JSON(404, gin.H{
			"status":  "failed",
			"message": "unable found items",
			"data":    err.Error(),
		})
		return
	}
	response := gin.H{
		"message":"Menu details fetched successfully",
		"menuID":menuItem.ID,
		"name":menuItem.Name,
		"category":menuItem.Category,
		"price":menuItem.Price,
		"duration":menuItem.Duration,
	}
	c.JSON(200, gin.H{"status":  "Success","menu":response})
}

// Create menulist for admin with authentication
func CreateMenu(c *gin.Context) {
	var menu models.MenuModel
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(menu)
	validate := validator.New()
	err := validate.Struct(menu)
	if err != nil {
		c.JSON(400, gin.H{
			"errror": "Binding error",
		})
		return
	}

	if err := database.DB.Create(&menu).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response := gin.H{
		"ItemID":menu.ID,
		"Category":menu.Category,
		"Price":menu.Price,
		"FoodImage":menu.FoodImage,
		"Duration":menu.Duration,
	}
	c.JSON(201, gin.H{"message": "Item added successfully", "data": response})
}

// Update the menu for admin with authentication
func UpdateMenu(c *gin.Context) {
	var menu models.MenuModel

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	menuID := c.Param("food_id")
	fmt.Println(menuID)
	var existingMenu models.MenuModel

	if err := database.DB.First(&existingMenu, menuID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not Found"})
		return
	}

	//update the fiels of existing menulist
	existingMenu.ID = menu.ID
	existingMenu.Category = menu.Category
	existingMenu.Name = menu.Name
	existingMenu.Price = menu.Price
	existingMenu.FoodImage = menu.FoodImage
	existingMenu.Duration = menu.Duration


	//save the updated menu item to the database

	if err := database.DB.Save(&existingMenu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update menu"})
		return
	}

	response := gin.H{
		"ItemID":menu.ID,
		"Category":menu.Category,
		"Price":menu.Price,
		"FoodImage":menu.FoodImage,
		"Duration":menu.Duration,
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Menu Details Updated successfully",
		"data":    response,
	})
}

// Delete the menu for admin with authentication
func DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	var menu models.MenuModel

	if err := database.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Failed",
			"message": "Menu id Not Found",
			"data":    err.Error(),
		})
		return
	}
	database.DB.Delete(&menu)
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Item Removed Successfully",
		"data":    id,
	})

}
