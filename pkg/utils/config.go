package utils

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the structure that holds the application configuration.
type Config struct {
	Database     DatabaseConfig `yaml:"database"`
	Kafka        KafkaConfig    `yaml:"kafka"`
	Redis        RedisConfig    `yaml:"redis"`
	Email        EmailConfig    `yaml:"email"`
	Jwt          JwtConfig      `yaml:"jwt"`
	RolesAllowed []string       `yaml:"roleAllowed"`
}

type JwtConfig struct {
	SecretKey string        `yaml:"secret_key"`
	Expiry    time.Duration `yaml:"expiry"`
}

type EmailConfig struct {
	SMTPHost     string `yaml:"smtpHost"`
	SMTPPort     string `yaml:"smtpPort"`
	SenderEmail  string `yaml:"senderEmail"`
	SenderPasswd string `yaml:"senderPasswd"`
}

// DatabaseConfig holds the database connection details.
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// KafkaConfig holds Kafka connection details.
type KafkaConfig struct {
	Brokers       []string `yaml:"brokers"`
	ConsumerGroup string   `yaml:"consumerGroup"`
	Topics        struct {
		Orders            string `yaml:"orders"`
		Inventory         string `yaml:"inventory"`
		Payments          string `yaml:"payments"`
		EmailConfirmation string `yaml:"emailConfirmation"`
	} `yaml:"topics"`
}

// RedisConfig holds Redis connection details.
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

// LoadConfig loads the YAML config file into a Config struct.
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Error decoding config file: %v", err)
		return nil, err
	}
	return &config, nil
}
