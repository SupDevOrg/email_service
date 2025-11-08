package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"email-service/internal/config"
	"email-service/internal/email"
	"email-service/internal/kafka"
	"email-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	emailService := email.NewEmailService(cfg)

	if err := emailService.TestConnection(); err != nil {
		log.Printf("Warning: SMTP connection test failed: %v", err)
	} else {
		log.Println("SMTP connection test: OK")
	}

	kafkaConsumer := kafka.NewConsumer(
		[]string{cfg.KafkaBrokers},
		cfg.KafkaTopic,          
		cfg.KafkaGroupID, 
		emailService,
	)
	emailHandler := handlers.NewEmailHandler(emailService)
	router := gin.Default()
	router.GET("/health", emailHandler.HealthCheck)
	go func() {
    router.Run(":" + cfg.HealthCheckPort) 
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafkaConsumer.Start(ctx)

	log.Printf("Email service started. Consuming from Kafka topic: %s", cfg.KafkaTopic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down email service...")
	kafkaConsumer.Close()
	log.Println("Email service stopped")
}
