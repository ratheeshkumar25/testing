package middleware

import (
	"fmt"
	"github.com/go-gomail/gomail"
	"io"
	"os"
)

// // SendEmail sends an email with an optional attachment
func SendEmail(msg, email, attachmentName string, attachmentData []byte) error {
	// SMTP server configuration
	senderEmail := os.Getenv("EMAIL")
	senderPassword := os.Getenv("PASSWORD")

	// Compose email message
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Confirmation-Email")
	m.SetBody("text/plain", msg)

	
	// m.Attach(attachmentName, gomail.SetCopyFunc(func(w gomail.) error {
	// 	_, err := w.Write(attachmentData)
	// 	return err
	// }))
    
	// Add attachment
	m.Attach(attachmentName, gomail.SetCopyFunc(func(w io.Writer) error {
		// _, err := buf.WriteTo(w)
		_, err := w.Write(attachmentData)
		return err
	}))

	// Dial to SMTP server and send email
	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPassword)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}


