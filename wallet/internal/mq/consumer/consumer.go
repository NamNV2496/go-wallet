package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/logic"
	"github.com/namnv2496/go-wallet/internal/mq"
)

type Consumer struct {
	transferService logic.TransferLogic
	consumer        sarama.ConsumerGroup
}

func NewConsumer(
	transferService logic.TransferLogic,
	config config.Config,
) (*Consumer, error) {

	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = mq.ClientId
	// saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = config.KafkaBroker
	}
	saramaConsumer, err := sarama.NewConsumerGroup([]string{kafkaBroker}, mq.ClientId, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		transferService: transferService,
		consumer:        saramaConsumer,
	}, nil
}

func (c *Consumer) Start() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	ctx := context.Background()

	go func() {
		for {
			if err := c.consumer.Consume(ctx, []string{mq.TopicResult}, c); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-signals
	log.Println("Interrupt is detected")
	close(doneCh)
	return nil
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Println("Received message:", string(msg.Key), ":", string(msg.Value))
		var transferResult mq.TransferResponse
		err := json.Unmarshal(msg.Value, &transferResult)
		if err != nil {
			log.Println("Failed to unmarshal Kafka message:", err)
			continue
		}
		ctx := context.Background()
		// check transferId
		_, err = c.transferService.GetTransfer(ctx, transferResult.TransferId)
		if err != nil {
			log.Println("TransferId is not exist: ", err)
			continue
		}

		if transferResult.Amount > 0 && transferResult.Status == 1 {
			if err := c.transferService.UpdateBalanceTransferTx(
				ctx,
				transferResult.FromId,
				transferResult.ToId,
				transferResult.Amount,
				transferResult.TransferId,
				transferResult.Status,
				transferResult.Message,
			); err != nil {
				log.Println("Failed to update balance:", err)
				// save to DLQ to handler later
				// continue
			}
		}
		// can trigger to app through notification or websocket

		// mark handler done message to stop retry
		session.MarkMessage(msg, "")
	}
	return nil
}
