package producer

import (
	"encoding/json"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/mq"
)

type Producer struct {
	client sarama.SyncProducer
}

func NewProducer(
	config config.Config,
) (*Producer, error) {

	configSarama := sarama.NewConfig()
	configSarama.Producer.RequiredAcks = sarama.WaitForAll
	configSarama.Producer.Retry.Max = mq.MaxRetry
	configSarama.Producer.Return.Successes = true
	configSarama.ClientID = mq.ClientId

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = config.KafkaBroker
	}
	producer, err := sarama.NewSyncProducer([]string{kafkaBroker}, configSarama)
	if err != nil {
		return nil, err
	}

	return &Producer{
		client: producer,
	}, nil
}

func (p *Producer) SendMessage(topic string, req any) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}
	partition, offset, err := p.client.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", mq.TopicRequest, partition, offset)
	return nil
}
