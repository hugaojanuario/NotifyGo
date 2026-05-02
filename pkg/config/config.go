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

type Config struct {
	Postgres PostgresConfig
	Kafka    KafkaConfig
	Api      ApiConfig
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
	}
}
