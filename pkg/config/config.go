package config

import (
	"os"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers string
	GroupID string
}

type ApiConfig struct {
	Secret string
	Port   string
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	From       string
}

type Config struct {
	Postgres PostgresConfig
	Kafka    KafkaConfig
	Api      ApiConfig
	SMTP     SMTPConfig
	Twilio   TwilioConfig
}

func LoadDotEnv() *Config {
	godotenv.Load("./.env")

	return &Config{
		Postgres: PostgresConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
		},

		Kafka: KafkaConfig{
			Brokers: os.Getenv("KAFKA_BROKERS"),
			GroupID: os.Getenv("KAFKA_GROUP_ID"),
		},

		Api: ApiConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Port:   os.Getenv("PORT"),
		},

		SMTP: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			User:     os.Getenv("SMTP_USER"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
		},

		Twilio: TwilioConfig{
			AccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
			AuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
			From:       os.Getenv("TWILIO_FROM"),
		},
	}
}
