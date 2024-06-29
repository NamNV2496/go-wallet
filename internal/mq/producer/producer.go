package producer

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/namnv2496/go-wallet/internal/mq"
)

type Producer struct {
	client sarama.SyncProducer
}

func NewProducer() (*Producer, error) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = mq.MaxRetry
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{mq.BrokerList}, config)
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
