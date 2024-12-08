package kafka

import (
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

// KafkaClient holds the producer and consumer clients.
type KafkaClient struct {
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

// InitKafka initializes Kafka producer and consumer.
func InitKafka(config *utils.KafkaConfig) (*KafkaClient, error) {
	// Initialize producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": strings.Join(config.Brokers, ",")})
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
		return nil, err
	}

	// Initialize consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Brokers, ","),
		"group.id":          "coffee-shop-consumer-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
		return nil, err
	}

	return &KafkaClient{Producer: producer, Consumer: consumer}, nil
}

// Close closes the Kafka producer and consumer.
func (k *KafkaClient) Close() {
	k.Producer.Close()
	k.Consumer.Close()
}
