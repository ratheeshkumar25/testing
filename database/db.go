package database

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"restaurant/models"
)

// type DBInterface interface {
//     Create(value interface{}) *gorm.DB
//     First(dest interface{}, conds ...interface{}) *gorm.DB
//     Save(value interface{}) *gorm.DB
//     Delete(value interface{}) *gorm.DB
//     Find(dest interface{}, conds ...interface{}) *gorm.DB
//     // Add other methods as needed
// }

// var DB DBInterface

var DB *gorm.DB

// DBConnect initializes the database connection.
func DBconnect() {
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	DB = db

	DB.AutoMigrate(
		&models.UsersModel{},
		&models.AdminModel{},
		&models.InvoicesModel{},
		&models.MenuModel{},
		&models.ReviewModel{},
		&models.StaffModel{},
		&models.TablesModel{},
		&models.RazorPay{},
		&models.ReservationModels{},
	)

}

// GetOrderByID retrieves an order by its ID.
func GetOrderByID(orderID uint) (*models.InvoicesModel, error) {
	var order models.InvoicesModel
	if err := DB.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetMenuByID retrieves a menu item by its ID.
func GetMenuByID(menuID uint) (*models.MenuModel, error) {
	var menu models.MenuModel
	if err := DB.First(&menu, menuID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

// GetUsersByID retrieves a user by their ID.
func GetUsersByID(userID uint) (*models.UsersModel, error) {
	// Return nil if userID is 0
	if userID == 0 {
		return nil, nil
	}

	var user models.UsersModel
	if err := DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetReservationByID retrieves a reservation by user ID.
func GetReservationByID(userID uint) (*models.ReservationModels, error) {
	var reservation models.ReservationModels
	if err := DB.Where("user_id = ?", userID).Find(&reservation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &reservation, nil
}
