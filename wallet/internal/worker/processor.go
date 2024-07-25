package worker

import (
	"context"

	"log"

	"github.com/hibiken/asynq"
	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type RedisTaskProcessor interface {
	StartQueue() error
	Shutdown()
	NewSendVerifyEmailTask(ctx context.Context, payload EmailVerifyPayload, opts ...asynq.Option) error
}

type TaskProcessor struct {
	client   *asynq.Client
	config   config.Config
	datapool *databaseaccess.Database
	server   *asynq.Server
}

func NewTaskProcessor(
	config config.Config,
	redisConfig *RedisConfigOpt,
	datapool *databaseaccess.Database,
) RedisTaskProcessor {

	client := asynq.NewClient(redisConfig.redisOpt)
	server := asynq.NewServer(
		redisConfig.redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Println("Error when create new server handler!")
			}),
			// Logger: logger,
		},
	)
	return &TaskProcessor{
		client:   client,
		config:   config,
		datapool: datapool,
		server:   server,
	}
}

func (t TaskProcessor) StartQueue() error {
	// start http server to trigger handler task
	mux := asynq.NewServeMux()

	// register handler request from http
	mux.HandleFunc(TaskSendVerifyEmail, t.HandlerSendEmailTask)

	return t.server.Start(mux)
}

func (t TaskProcessor) Shutdown() {

	t.server.Shutdown()
	defer t.client.Close()
}
