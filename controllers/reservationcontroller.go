package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"restaurant/database"
	"restaurant/middleware"
	"restaurant/models"

	//"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

func CreateReservartion(c *gin.Context) {
	var reservation models.ReservationModels
	if err := c.BindJSON(&reservation); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Check if the user already has a reservation for the requested time slot
	userIDContext, _ := c.Get("userID")
	userID := userIDContext.(uint)

	var existingReservation models.ReservationModels
	err := database.DB.Where("user_id = ? AND ((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?) OR (end_time > ? AND end_time <= ?))", userID, reservation.EndTime, reservation.StartTime, reservation.StartTime, reservation.EndTime, reservation.StartTime, reservation.EndTime).First(&existingReservation).Error
	if err == nil {
		c.JSON(400, gin.H{"error": "You already have a reservation for the requested time slot"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(500, gin.H{"error": "Failed to check existing reservation"})
		return
	}

	// Check for available tables and staff
	availableTable, availableStaff, err := checkAvailability(reservation.StartTime, reservation.EndTime, reservation.NumberOfGuest)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Assign the available table and staff to the reservation
	reservation.TableID = uint(availableTable.ID)
	reservation.StaffID = availableStaff.ID
	reservation.UserID = userID

	reservation.UserID = userID
	if err := database.DB.Create(&reservation).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//Generate PDF for the reservation
	pdfBytes, err := GeneratePDFReservation(reservation)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate the Reservartion PDF"})
		return
	}

	//Send the reservation confirmation email with PDF invoice attached
	err = middleware.SendEmail("Reservation done successfully", reservation.Email, "reservation.pdf", pdfBytes)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(201, gin.H{
		"message":       "Reservation created successfully-Confirmation email send",
		"reservation":   reservation.ID,
		"customerID":    reservation.UserID,
		"selectedGuest": reservation.NumberOfGuest,
		"table":         reservation.TableID,
		"startTime":     reservation.StartTime,
		"endTime":       reservation.EndTime,
		"staffServe":    reservation.StaffID,
	})

}

func UpdateReservation(c *gin.Context) {
	//Get the reservation ID from the request URL
	reservationID := c.Param("id")
	//Reterive the existing reservation from the database
	var reservation models.ReservationModels
	if err := database.DB.First(&reservation, reservationID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Reservation not found"})
		return
	}
	//Parase updated reservation data  from the requestbody
	var updatedReservation models.ReservationModels
	if err := c.BindJSON(&updatedReservation); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check for available tables and staff members for the new time slot
	availableTable, availableStaff, err := checkAvailability(updatedReservation.StartTime, updatedReservation.EndTime, reservation.NumberOfGuest)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update the reservation with the new table and staff assignments
	reservation.TableID = uint(availableTable.ID)
	reservation.StaffID = availableStaff.ID
	reservation.StartTime = updatedReservation.StartTime
	reservation.EndTime = updatedReservation.EndTime

	//Get userID from the context
	userIDContext, _ := c.Get("userID")
	fmt.Println(userIDContext)
	userID := userIDContext.(uint)

	var bookingID models.UsersModel

	if err := database.DB.Where("user_id = ?", userID).First(&bookingID).Error; err == gorm.ErrRecordNotFound {
		c.JSON(200, gin.H{"status": "Success",
			"message": "No booking",
			"data":    nil})
		return
	} else if err != nil {
		c.JSON(404, gin.H{"status": "Failed",
			"message": "Database error",
			"data":    nil})
		return
	}
	reservation.UserID = userID
	if err := database.DB.Save(&reservation).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Generate PDF for the updated reservation
	pdfBytes, err := GeneratePDFReservation(reservation)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate PDF"})
		return
	}

	// Send email notification with the updated reservation PDF attached
	err = middleware.SendEmail("Reservation updated successfully", reservation.Email, "invoice.pdf", pdfBytes)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send email"})
		return
	}

	// Respond with updated reservation details
	c.JSON(200, gin.H{
		"message":       "Reservation updated successfully",
		"reservation":   reservation.ID,
		"customerID":    reservation.UserID,
		"selectedGuest": reservation.NumberOfGuest,
		"table":         reservation.TableID,
		"startTime":     reservation.StartTime,
		"endTime":       reservation.EndTime,
		"staffServe":    reservation.StaffID,
	})
}

func CancelReservation(c *gin.Context) {
	reservationID := c.Param("id")
	var reservation models.ReservationModels

	if err := database.DB.First(&reservation, reservationID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Reservation not found"})
		return
	}

	if err := database.DB.Delete(&reservation).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Reservation canceled successfully",
	})
}
func SearchAvailableTables(c *gin.Context) {
	var request struct {
		StartTime     time.Time `json:"startTime"`
		EndTime       time.Time `json:"endTime"`
		NumberOfGuest int       `json:"numberOfGuest"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Fetch available tables based on guest count
	var availableTables []models.TablesModel
	err := database.DB.Where("capacity >= ? AND availability = true", request.NumberOfGuest).Find(&availableTables).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Fetch reservations that overlap with the requested time slot
	var overlappingReservations []models.ReservationModels
	err = database.DB.Where("(start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?) OR (end_time > ? AND end_time <= ?)", request.EndTime, request.StartTime, request.StartTime, request.EndTime, request.StartTime, request.EndTime).Find(&overlappingReservations).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Separate available and booked tables
	var availableTableIDs, bookedTableIDs []uint
	bookedTables := make(map[uint][]models.ReservationModels)

	for _, reservation := range overlappingReservations {
		bookedTableIDs = append(bookedTableIDs, reservation.TableID)
		bookedTables[reservation.TableID] = append(bookedTables[reservation.TableID], reservation)
	}

	for _, table := range availableTables {
		if containsUint(bookedTableIDs, table.ID) {
			continue
		}
		availableTableIDs = append(availableTableIDs, table.ID)
	}

	// Prepare response data for available tables
	var availableTableData []gin.H
	for _, tableID := range availableTableIDs {
		table, err := getTableCapacityByID(tableID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		availableTableData = append(availableTableData, gin.H{
			"table_id":     table.ID,
			"capacity":     table.Capacity,
			"availability": true,
		})
	}

	// Prepare response data for booked tables with reservation details
	var bookedTableData []gin.H
	for _, reservation := range overlappingReservations {
		bookedTableData = append(bookedTableData, gin.H{
			"table_id":      reservation.TableID,
			"reservationID": reservation.ID,
			"numberofGuest": reservation.NumberOfGuest,
			"startTime":     reservation.StartTime,
			"endTime":       reservation.EndTime,
		})
	}

	c.JSON(200, gin.H{
		"message":          "Tables availability information fetched successfully",
		"available_tables": availableTableData,
		"booked_tables":    bookedTableData,
	})
}

// Custom function to check if a slice of uint contains a specific uint value
func containsUint(slice []uint, val uint) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func getTableCapacityByID(tableID uint) (*models.TablesModel, error) {
	var table models.TablesModel
	if err := database.DB.Where("id = ?", tableID).Find(&table).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &table, nil
}

func checkAvailability(startTime, endTime time.Time, numGuests int) (*models.TablesModel, *models.StaffModel, error) {
	var availableTable models.TablesModel
	var availableStaff models.StaffModel

	// Check for available tables with sufficient capacity
	err := database.DB.Where("capacity >= ? AND availability = true", numGuests).First(&availableTable).Error
	if err != nil {
		return nil, nil, fmt.Errorf("no available tables found")
	}

	// Get all existing reservations for the requested time slot
	var existingReservations []models.ReservationModels
	err = database.DB.Where("(start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?) OR (end_time > ? AND end_time <= ?)", endTime, startTime, startTime, endTime, startTime, endTime).Find(&existingReservations).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, fmt.Errorf("failed to retrieve existing reservations: %v", err)
	}

	// Check if the selected table is not already reserved for the requested duration
	for _, reservation := range existingReservations {
		if reservation.TableID == availableTable.ID {
			return nil, nil, fmt.Errorf("selected table is not available for the requested time slot")
		}
	}

	// Get all available staff members
	var availableStaffMembers []models.StaffModel
	err = database.DB.Where("blocked = false").Find(&availableStaffMembers).Error
	if err != nil {
		return nil, nil, fmt.Errorf("no available staff found")
	}

	// Check if a staff member is available for the requested duration
	for _, staff := range availableStaffMembers {
		var reservationsForStaff []models.ReservationModels
		err = database.DB.Where("staff_id = ? AND ((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?) OR (end_time > ? AND end_time <= ?))", staff.ID, startTime, endTime, startTime, endTime, startTime, endTime).Find(&reservationsForStaff).Error
		if err != nil {
			return nil, nil, fmt.Errorf("failed to check staff availability: %v", err)
		}

		if len(reservationsForStaff) == 0 {
			availableStaff = staff
			break
		}
	}

	if availableStaff.ID == 0 {
		return nil, nil, fmt.Errorf("no staff available for the requested time slot")
	}

	return &availableTable, &availableStaff, nil
}

func GeneratePDFReservation(reservation models.ReservationModels) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 10)

	// Set text color to purple for the heading
	pdf.SetTextColor(128, 0, 128) // Purple color

	// Add "Go Restaurant" as header with hotel address
	pdf.CellFormat(0, 10, "Go Restaurant", "", 0, "L", false, 0, "")
	pdf.Ln(5) // Move down for spacing
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 10, "HSR, Layout, Bengaluru, Karnataka-560102", "", 0, "L", false, 0, "")
	pdf.Ln(10) // Move down for spacing

	// Add title "Reservation Confirmation"
	pdf.SetFont("Arial", "B", 10) // Set font size for the title
	pdf.CellFormat(0, 10, "Reservation Confirmation", "", 0, "C", false, 0, "")
	pdf.Ln(10) // Move down for spacing

	// Add subject line with big font
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 10, "Dear valued customer,", "", 0, "L", false, 0, "")
	pdf.Ln(3)
	pdf.CellFormat(0, 10, "Your table has been reserved and below are your reservation details:", "", 0, "L", false, 0, "")
	pdf.Ln(10) // Move down for spacing

	// Create a separate box for reservation details
	pdf.SetFillColor(220, 220, 220) // Light gray background color
	pdf.Rect(10, pdf.GetY(), 190, 60, "F")
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0) // Black text color

	// Add reservation details inside the box
	pdf.Ln(2) // Move down for spacing within the box
	pdf.CellFormat(0, 7, fmt.Sprintf("Reservation: %d", reservation.ID), "", 0, "L", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(0, 7, fmt.Sprintf("CustometrID: %d", reservation.UserID), "", 0, "L", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(0, 7, fmt.Sprintf("Table: %d", reservation.TableID), "", 0, "L", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(0, 7, fmt.Sprintf("Number of Guests: %d", reservation.NumberOfGuest), "", 0, "L", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(0, 7, fmt.Sprintf("Staff ID to Serve: %d", reservation.StaffID), "", 0, "L", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(0, 7, fmt.Sprintf("Start Time: %s", reservation.StartTime.Format("2006-01-02 15:04:05")), "", 0, "L", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(0, 7, fmt.Sprintf("End Time: %s", reservation.EndTime.Format("2006-01-02 15:04:05")), "", 0, "L", false, 0, "")

	// Add "Thank you for choosing. Welcome back again!" message
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(255, 0, 0) // Red text color
	pdf.Ln(10)                  // Move down for spacing
	pdf.CellFormat(0, 10, "Thank you for choosing. Welcome back again!", "", 0, "C", false, 0, "")
	pdf.Ln(5) // Move down for spacing

	// Add system-generated confirmation message
	pdf.CellFormat(0, 10, "This is a system-generated reservation confirmation.", "", 0, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
