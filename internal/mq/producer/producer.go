package producer

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/namnv2496/go-wallet/internal/mq"
)

func NewProducer() (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = mq.MaxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{mq.BrokerList}, config)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln("Failed to close producer")
		}
	}()
	msg := &sarama.ProducerMessage{
		Topic: mq.Topic,
		Value: sarama.StringEncoder("Something Cool"),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return nil, err
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", mq.Topic, partition, offset)
	return producer, nil
}
