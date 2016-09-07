package core

import (
	"errors"
	"github.com/Shopify/sarama"
	"time"
)

const (
	kafkaLogsTopic = "amp-logs"
)

var (
	kafkaClient sarama.Client
)

// Kafka singleton
type Kafka struct {
}

// Connect to kafka
func (kafka *Kafka) Connect(host string) error {
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_0_0

	var err error
	kafkaClient, err = sarama.NewClient([]string{host}, config)
	return err
}

// NewConsumer create a new consumer
func (kafka *Kafka) NewConsumer() (sarama.Consumer, error) {
	return sarama.NewConsumerFromClient(kafkaClient)
}

// NewAsyncProducer create a new async producer
func (kafka *Kafka) NewAsyncProducer() (sarama.AsyncProducer, error) {
	return sarama.NewAsyncProducerFromClient(kafkaClient)
}

// Topics return available topics
func (kafka *Kafka) Topics() ([]string, error) {
	kafkaClient.RefreshMetadata()
	return kafkaClient.Topics()
}

// WaitForTopic wait for given topic availability
func (kafka *Kafka) WaitForTopic(topic string, timeout int) error {
	topicFound := false
WaitForTopic:
	for i := 0; i < timeout; i++ {
		topics, err := kafka.Topics()
		if err != nil {
			return err
		}
		for _, topic := range topics {
			if topic == topic {
				topicFound = true
				break WaitForTopic
			}
		}
		time.Sleep(1 * time.Second)
	}

	if !topicFound {
		return errors.New("Kafka topic not available.")
	}
	return nil
}

// Close close the connection to Kafka
func (kafka *Kafka) Close() error {
	return kafkaClient.Close()
}
