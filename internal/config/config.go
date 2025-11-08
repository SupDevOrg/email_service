package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort   string
	HealthCheckPort  string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string

	KafkaBrokers string
	KafkaTopic   string
	KafkaGroupID string
}

func Load() *Config {
	return &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		HealthCheckPort:   getEnv("HEALTH_PORT", "8081"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.mail.ru"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", "supdev@list.ru"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "HHrq5jHranb5OXDhyYmy"),
		FromEmail:    getEnv("FROM_EMAIL", ""),

		KafkaBrokers: getEnv("KAFKA_BROKERS", "kafka:29092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "email-auth-codes"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "email-service"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
