package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"restaurant/database"
	"restaurant/middleware"
	"restaurant/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"gorm.io/gorm"
)

// User Login
func GetHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to RERA Restaurant World Please log in with your mobile"})
}

// Postloginhandler handles the login request
func PostLogin(c *gin.Context) {
	var users models.UsersModel
	if err := c.BindJSON(&users); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if user exists in database
	if err := database.DB.Where("phone = ?", users.Phone).First(&users).Error; err == nil {
		// // Generate Token
		token, err := middleware.GenerateUsertoken(users.Phone, uint(users.UserID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Token Generated Succesfully", "token": token})

		return
	} else if err != gorm.ErrRecordNotFound {

		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error exist"})
		return
	}
	// Send OTP via Twilio
	err := SendOTP(users.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP", "data": err.Error()})
		return
	}
	//Marshal json data
	userData, err := json.Marshal(&users)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal user", "data": err.Error()})
		return
	}

	key := fmt.Sprintf("user:%s", users.Phone)
	err = database.SetRedis(key, userData, time.Minute*5)
	if err != nil {
		fmt.Println("Error srtting user in Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Data": nil, "Message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP generated successfull go to verification page"})

}

// SendOTP is send the OTP via Twilio SMS
func SendOTP(phoneNumber string) error {
	//Load Twilio credentials from enviornment variable
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	//Create SMS message
	from := os.Getenv("TWILIO_PHONE_NUMBER")
	params := verify.CreateVerificationParams{}
	params.SetTo("+919961429910")
	params.SetChannel("sms")
	println(from)
	response, err := client.VerifyV2.CreateVerification(os.Getenv("SERVICE_TOKEN"), &params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(response)

	return nil
}

// VerifyOtp verifies the OTP sent to the users phone
func SignupVerify(c *gin.Context) {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	// request body data
	var verifyModel models.VerifyOTP
	if err := c.BindJSON(&verifyModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Success": false, "Data": nil, "Message": err.Error()})
		return
	}

	//checking for the missing otp request

	if verifyModel.Otp == "" {
		c.JSON(http.StatusOK, gin.H{"Success": true, "Message": "OTP verified successfully"})

	}
	// Create a Twilio REST client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	params := verify.CreateVerificationCheckParams{}
	params.SetTo("+919961429910")
	params.SetCode(verifyModel.Otp)

	//send twilio verification check
	resp, err := client.VerifyV2.CreateVerificationCheck(os.Getenv("SERVICE_TOKEN"), &params)
	if err != nil {
		fmt.Println("err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Data": nil, "Message": "error in verifying OTP provided"})
		return
	} else if *resp.Status != "approved" {
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Data": nil, "Message": "wrong OTP provided"})
		return
	}
	// fmt.Println("Twilio verification response", response.Status)

	// Check if user already exists in Redis
	key := fmt.Sprintf("user:%s", verifyModel.Phone)
	value, err := database.GetRedis(key)
	if err != nil {
		fmt.Println("Error checking user in Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Success": false, "Data": nil, "Message": "Internal server error"})
		return
	}
	// Bind json data to unmarshal
	var user models.UsersModel
	err = json.Unmarshal([]byte(value), &user)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal user", "data": err.Error()})
		return
	}

	err = database.DB.Create(&user).Error
	if err != nil {
		fmt.Println("Error creating user", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Status": false, "Data": nil, "Message": "Failed to create user"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"Status": true, "Message": "OTP verified successfully"})

}

// User logout
func UserLogout(c *gin.Context) {

	//After Successful LOGOUT
	c.JSON(http.StatusOK, gin.H{"message": "Logout successfully"})
}

