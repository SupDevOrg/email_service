package kafka

import (
	"context"
	"encoding/json"
	"log"

	"email-service/internal/email"

	"github.com/segmentio/kafka-go"
)

type AuthCodeMessage struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	Type  string `json:"type,omitempty"`
}

type Consumer struct {
	reader       *kafka.Reader
	emailService *email.EmailService
}

func NewConsumer(brokers []string, topic string, groupID string, emailService *email.EmailService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{
		reader:       reader,
		emailService: emailService,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("Starting Kafka consumer for topic: %s", c.reader.Config().Topic)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			go c.processMessage(msg)
		}
	}
}

func (c *Consumer) processMessage(msg kafka.Message) {
	var authMsg AuthCodeMessage
	if err := json.Unmarshal(msg.Value, &authMsg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	if authMsg.Type == "" {
		authMsg.Type = "login"
	}

	log.Printf("Received auth code request for email: %s, type: %s", authMsg.Email, authMsg.Type)

	if err := c.emailService.SendAuthCode(authMsg.Email, authMsg.Code, authMsg.Type); err != nil {
		log.Printf("Failed to send auth code to %s: %v", authMsg.Email, err)
	} else {
		log.Printf("Auth code sent successfully to %s", authMsg.Email)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
