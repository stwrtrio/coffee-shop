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
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Brokers, ","),
	})
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
		return nil, err
	}

	// Initialize consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Brokers, ","),
		"group.id":          config.ConsumerGroup,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
		return nil, err
	}

	return &KafkaClient{Producer: producer, Consumer: consumer}, nil
}

// Produce sends a message to a Kafka topic.
func (k *KafkaClient) Produce(topic string, message []byte) error {
	return k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
}

// Consume listens for messages on a Kafka topic and processes them using the provided handler.
func (k *KafkaClient) Consume(topic string, handler func(message *kafka.Message) error) {
	// Subscribe to the topic
	err := k.Consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	for {
		msg, err := k.Consumer.ReadMessage(-1) // -1: Wait indefinitely
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}

		// Handle the message
		if err := handler(msg); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	}
}

// Close closes the Kafka producer and consumer.
func (k *KafkaClient) Close() {
	k.Producer.Close()
	k.Consumer.Close()
}
