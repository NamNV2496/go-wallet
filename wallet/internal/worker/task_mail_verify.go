package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/namnv2496/go-wallet/config"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
	"github.com/namnv2496/go-wallet/internal/email"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type EmailVerifyPayload struct {
	Username   string `json:"username"`
	SecretCode string `json:"secret_code"`
	Email      string `json:"email"`
}

func (t TaskProcessor) NewSendVerifyEmailTask(
	ctx context.Context,
	payload EmailVerifyPayload,
	opts ...asynq.Option,
) error {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload)
	// enqueue to execute immediate
	info, err := t.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Println("Register new task type: ", task.Type(), "queue: ", info.Queue)
	return nil
}

func (t TaskProcessor) HandlerSendEmailTask(
	ctx context.Context,
	task *asynq.Task,
) error {

	log.Println("HandlerSendEmailTask was called")
	// create CreateVerifyEmail record
	var payload EmailVerifyPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	arg := db.CreateVerifyEmailParams{
		Username:   payload.Username,
		Email:      payload.Email,
		SecretCode: payload.SecretCode,
	}
	_, err := t.datapool.CreateVerifyEmail(ctx, arg)
	if err != nil {
		return err
	}

	// send email
	sendEmail(t.config, payload.Username, payload.Email, payload.SecretCode)
	log.Println("Send email")
	return nil
}

func sendEmail(config config.Config, username string, toEmail string, secretCode string) error {
	mailer := email.NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword,
	)
	subject := "Welcome to Go-wallet"
	verifyUrl := fmt.Sprintf("http://localhost:8080/api/v1/verify_email?username=%s&secret_code=%s", username, secretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>`, username, verifyUrl)
	to := []string{toEmail}

	err := mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}
	return nil
}
