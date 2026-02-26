package config

import (
	"os"
	"strconv"
	"fmt"
)

type Config struct {
	ServerPort   string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string

	KafkaBrokers string
	KafkaTopic   string
	KafkaGroupID string
}

func Load() (*Config, error) {
	cfg := &Config{}

	var ok bool
	if cfg.ServerPort, ok = getEnv("SERVER_PORT"); !ok {
		return nil, fmt.Errorf("SERVER_PORT is required")
	}

	if cfg.SMTPHost, ok = getEnv("SMTP_HOST"); !ok {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}

	if cfg.SMTPUsername, ok = getEnv("SMTP_USERNAME"); !ok {
		return nil, fmt.Errorf("SMTP_USERNAME is required")
	}

	if cfg.SMTPPassword, ok = getEnv("SMTP_PASSWORD"); !ok {
		return nil, fmt.Errorf("SMTP_PASSWORD is required")
	}

	if cfg.FromEmail, ok = getEnv("FROM_EMAIL"); !ok {
		return nil, fmt.Errorf("FROM_EMAIL is required")
	}

	if cfg.KafkaBrokers, ok = getEnv("KAFKA_BROKERS"); !ok {
		return nil, fmt.Errorf("KAFKA_BROKERS is required")
	}

	if cfg.KafkaTopic, ok = getEnv("KAFKA_TOPIC"); !ok {
		return nil, fmt.Errorf("KAFKA_TOPIC is required")
	}

	if cfg.KafkaGroupID, ok = getEnv("KAFKA_GROUP_ID"); !ok {
		return nil, fmt.Errorf("KAFKA_GROUP_ID is required")
	}

	cfg.SMTPPort = getEnvAsInt("SMTP_PORT", 0)
	if cfg.SMTPPort == 0 {
		return nil, fmt.Errorf("SMTP_PORT is required and must be int")
	}

	return cfg, nil
}

func getEnv(key string) (string, bool) {
	value, exists := os.LookupEnv(key)
	return value, exists
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
