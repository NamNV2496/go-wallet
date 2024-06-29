package consumer

import (
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/namnv2496/go-wallet/internal/mq"
)

func NewConsumer() error {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := mq.BrokerList
	master, err := sarama.NewConsumer([]string{brokers}, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Fatalln("Failed to close consumer")
		}
	}()
	consumer, err := master.ConsumePartition(mq.Topic, int32(mq.Partition), sarama.OffsetOldest)
	if err != nil {
		return err
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				log.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed messages")
	return nil
}
