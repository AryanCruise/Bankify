package notifications

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// NotificationService defines the service responsible for handling notifications.
type NotificationService struct {
	Email             string
	Password          string
	Host              string
	SMTPPort          int
	TwilioSID         string
	TwilioAuthToken   string
	TwilioPhoneNumber string
}

// NewNotificationService creates a new NotificationService instance.
func NewNotificationService(email, password, host string, port int, twilioSID, twilioAuthToken, twilioPhoneNumber string) *NotificationService {
	return &NotificationService{
		Email:             email,
		Password:          password,
		Host:              host,
		SMTPPort:          port,
		TwilioSID:         twilioSID,
		TwilioAuthToken:   twilioAuthToken,
		TwilioPhoneNumber: twilioPhoneNumber,
	}
}

// SendNotification sends the notification based on the message content.
func (s *NotificationService) SendEmailNotification(recieverEmail, message string) {
	// Extend this to include email/SMS sending logic
	from := s.Email
	to := []string{recieverEmail}
	subject := "New Transaction Notification"
	body := fmt.Sprintf("Subject: %s \n\n%s", subject, message)

	auth := smtp.PlainAuth("", from, s.Password, s.Host)

	serverAddr := fmt.Sprintf("%s:%d", s.Host, s.SMTPPort)
	fmt.Println("SMTP Server Address:", serverAddr)

	err := smtp.SendMail(
		serverAddr,
		auth,
		from,
		to,
		[]byte(body),
	)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return
	}
	fmt.Printf("Notification sent to %s\n successfully", to)
}

func (s *NotificationService) SendSMSNotification(recipientCountryCode, recipientPhoneNumber, message string) {
	// Initialize Twilio client

	if !strings.HasPrefix(recipientPhoneNumber, "+") {
		recipientPhoneNumber = recipientCountryCode + recipientPhoneNumber
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: s.TwilioSID,
		Password: s.TwilioAuthToken,
	})

	// Create message parameters
	params := &openapi.CreateMessageParams{
		From: &s.TwilioPhoneNumber,  // Your Twilio phone number
		To:   &recipientPhoneNumber, // Recipient's phone number
		Body: &message,              // Message content
	}

	// Send the SMS message using Twilio's API
	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
		return
	}

	// Output the response from Twilio
	fmt.Printf("SMS Notification sent to %s successfully: SID=%s\n", recipientPhoneNumber, *resp.Sid)
}
