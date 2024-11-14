package notifications

import (
	database "Accounts/internal/datab"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type TransactionMessage struct {
	UserID  int32  `json:"user_id"`
	Message string `json:"message"`
}

func formatNotificationMessage(txnMsg TransactionMessage) string {
	return fmt.Sprintf(
		"User ID: %d\nTransaction Details:\n\n%s",
		txnMsg.UserID,
		txnMsg.Message,
	)
}

// StartConsumer initializes a Kafka consumer to listen to a specific topic.
func StartConsumer(brokerList []string, topic string, notificationService NotificationService, db *database.Queries, ctx context.Context) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}

	defer partitionConsumer.Close()

	fmt.Printf("Consuming messages from topic: %s\n", topic) // Log topic

	for message := range partitionConsumer.Messages() {
		fmt.Printf("Received message: %s\n", string(message.Value))
		var txnMsg TransactionMessage
		err := json.Unmarshal(message.Value, &txnMsg)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Retrieve the email from the database based on the userID
		UserInfo, err := db.GetUserByUserID(ctx, txnMsg.UserID)
		if err != nil {
			log.Printf("Failed to get user email: %v", err)
			continue
		}

		formattedMessage := formatNotificationMessage(txnMsg)
		notificationService.SendEmailNotification(UserInfo.EmailID, formattedMessage)
		notificationService.SendSMSNotification(UserInfo.CountryCode, UserInfo.PhoneNo, formattedMessage)
		fmt.Println("Notification sent successfully!")
	}
}
