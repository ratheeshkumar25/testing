package models

import "time"

// Invoice model represents an invoice entity.
type InvoicesModel struct {
	InvoiceID      int       `gorm:"primaryKey;autoIncrement"`
	OrderID        int       `gorm:"autoIncrement"`
	TableID        int       `json:"tableID"`
	StaffID        int       `json:"staffID"`
	Quantity       int       `json:"quantity"`
	Email          string    `json:"email"`
	TotalAmount    float64   `json:"totalAmount"`
	PaymentMethod  string    `json:"paymentMethod"`
	PaymentDueDate time.Time `json:"paymentDueDate"`
	PaymentStatus  string    `json:"paymentStatus"`
	ItemID         uint
	UserID         uint
}

// RazorPay model represents RazorPay payment details.
type RazorPay struct {
	InvoiceID       uint    `JSON:"userID"`
	RazorPaymentID  string  `JSON:"razorpaymentID" gorm:"primaryKey;autoIncrement"`
	RazorPayOrderID string  `JSON:"razorpayorderID"`
	Signature       string  `JSON:"signature"`
	AmountPaid      float64 `JSON:"amountpaid"`
}
