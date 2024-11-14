package notifications

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// InitializeProducer initializes a new Kafka producer.
func InitializeProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Printf("Failed to initialize producer: %v", err)
		return nil, err
	}
	return producer, nil
}

// SendNotification sends a transaction notification to the Kafka topic.
func SendNotification(producer sarama.SyncProducer, userID int32, topic, message string) error {
	// Check if producer is nil
	if producer == nil {
		log.Println("Error: Kafka producer is not initialized")
		return fmt.Errorf("kafka producer is not initialized")
	}

	notificationMessage := struct {
		UserID  int32  `json:"user_id"`
		Message string `json:"message"`
	}{
		UserID:  userID,
		Message: message,
	}

	// Serialize the message to JSON
	jsonMessage, err := json.Marshal(notificationMessage)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMessage),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	fmt.Printf("Message sent to topic %s: %s\n", topic, message)
	return nil
}
