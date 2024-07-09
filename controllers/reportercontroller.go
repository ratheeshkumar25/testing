package controllers

import (
	"math"
	"restaurant/database"
	"restaurant/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get the the total order
func TotalOrder(c *gin.Context) {
	//Fetch the order details count from database
	var totalOrder int64
	if err := database.DB.Model(&models.InvoicesModel{}).Count(&totalOrder).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch the order"})
		return
	}
	c.JSON(200, gin.H{"totalOrder": totalOrder})
}

// Get total sales with additional details
func TotalSales(c *gin.Context) {
	type SalesData struct {
		TotalSales           string `json:"totalSales"`
		CompletedOrders      int64  `json:"completedOrders"`
		PendingOrders        int64  `json:"pendingOrders"`
		AvgOrderValue        string `json:"avgOrderValue"`
		AvgPendingOrderValue string `json:"avgPendingOrderValue"`
	}

	var salesData SalesData
	rows, err := database.DB.Model(&models.InvoicesModel{}).
		Select(`
            SUM(CASE WHEN payment_status = 'Completed' THEN total_amount ELSE 0 END) as total_sales,
            COUNT(CASE WHEN payment_status = 'Completed' THEN 1 END) as completed_orders,
            COUNT(CASE WHEN payment_status = 'Pending' THEN 1 END) as pending_orders,
            SUM(CASE WHEN payment_status = 'Pending' THEN total_amount ELSE 0 END) as pending_sales
        `).Rows()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch the sales report"})
		return
	}
	defer rows.Close()

	if rows.Next() {
		var totalSales, pendingSales float64
		if err := rows.Scan(&totalSales, &salesData.CompletedOrders, &salesData.PendingOrders, &pendingSales); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan the sales report"})
			return
		}
		salesData.TotalSales = strconv.FormatFloat(totalSales, 'f', 2, 64)
		if salesData.CompletedOrders > 0 {
			avgOrderValue := totalSales / float64(salesData.CompletedOrders)
			salesData.AvgOrderValue = strconv.FormatFloat(avgOrderValue, 'f', 2, 64)
		}
		if salesData.PendingOrders > 0 {
			avgPendingOrderValue := pendingSales / float64(salesData.PendingOrders)
			salesData.AvgPendingOrderValue = strconv.FormatFloat(avgPendingOrderValue, 'f', 2, 64)
		}
	}

	c.JSON(200, salesData)
}

// Employee perfomance
func EmployeePerfomance(c *gin.Context) {
	var employeePerfomance []struct {
		StaffID     uint `json:"staffId"`
		TotalOrders int  `json:"totalOrders"`
	}

	if err := database.DB.Model(&models.InvoicesModel{}).Select("staff_id, COUNT(*) as total_orders").Group("staff_id").Scan(&employeePerfomance).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch employee performance"})
		return
	}

	c.JSON(200, gin.H{"employeePerformance": employeePerfomance})
}

// Get the most ordered items with menu IDs
func MostOrderedItems(c *gin.Context) {
	type OrderedItem struct {
		ItemID   uint `json:"itemId"`
		Quantity int  `json:"total_quantity"`
	}

	var orderedItems []OrderedItem
	if err := database.DB.Model(&models.InvoicesModel{}).
		Select("item_id, SUM(quantity) as total_quantity").
		Group("item_id").
		Order("total_quantity DESC").
		Limit(10).
		Scan(&orderedItems).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch most ordered items"})
		return
	}
	c.JSON(200, gin.H{"mostOrderedItems": orderedItems})
}

func RevenueReport(c *gin.Context) {
	type RevenueData struct {
		TotalRevenue           float64 `json:"totalRevenue"`
		TotalOrders            int64   `json:"totalOrder"`
		AverageRevenuePerOrder float64 `json:"averageRevenueperOrder"`
		TotalSalaryExpense     float64 `json:"totalSalaryExpense"`
		Profit                 float64 `json:"profit"`
	}
	var revenueData RevenueData

	//Calculate the total Revenue
	var totalRevenue float64
	if err := database.DB.Model(&models.InvoicesModel{}).
		Select("SUM(total_amount)").
		Row().
		Scan(&totalRevenue); err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch total revenue"})
		return
	}
	revenueData.TotalRevenue = totalRevenue

	// Calculate total orders
	var totalOrders int64
	if err := database.DB.Model(&models.InvoicesModel{}).
		Count(&totalOrders).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch total orders"})
		return
	}
	revenueData.TotalOrders = totalOrders

	// Calculate average revenue per data
	if totalOrders > 0 {
		revenueData.AverageRevenuePerOrder = totalRevenue / float64(totalOrders)
	}
	// Calculate total salary expense
	var totalSalaryExpense float64
	if err := database.DB.Model(&models.StaffModel{}).
		Select("SUM(salary)").
		Row().
		Scan(&totalSalaryExpense); err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch total salary expense"})
		return
	}
	revenueData.TotalSalaryExpense = totalSalaryExpense

	//calculate total profit
	revenueData.Profit = totalRevenue - totalSalaryExpense

	// Format the decimal values to two decimal points
	revenueData.TotalRevenue = roundToTwoDecimalPlaces(revenueData.TotalRevenue)
	revenueData.AverageRevenuePerOrder = roundToTwoDecimalPlaces(revenueData.AverageRevenuePerOrder)
	revenueData.Profit = roundToTwoDecimalPlaces(revenueData.Profit)

	c.JSON(200, revenueData)
}

// roundToTwoDecimalPlaces rounds a float64 value to two decimal places
func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}
