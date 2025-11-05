package provider

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tnqbao/gau-account-service/shared/infra"
)

type EmailMessage struct {
	Type          string `json:"type"`
	Recipient     string `json:"recipient"`
	RecipientName string `json:"recipientName,omitempty"`
	Content       string `json:"content"`
	ActionUrl     string `json:"actionUrl,omitempty"`
}

type EmailProducer struct {
	rabbitmq *infra.RabbitMQClient
}

func NewEmailProducer(rabbitmq *infra.RabbitMQClient) *EmailProducer {
	return &EmailProducer{
		rabbitmq: rabbitmq,
	}
}

func (p *EmailProducer) SendEmailConfirmation(ctx context.Context, email, recipientName, content, actionUrl string) error {
	message := EmailMessage{
		Type:          "confirmation",
		Recipient:     email,
		RecipientName: recipientName,
		Content:       content,
		ActionUrl:     actionUrl,
	}

	return p.publishEmail(ctx, "email.confirmation", message)
}

func (p *EmailProducer) SendEmailNotification(ctx context.Context, email, recipientName, content, actionUrl string) error {
	message := EmailMessage{
		Type:          "notification",
		Recipient:     email,
		RecipientName: recipientName,
		Content:       content,
		ActionUrl:     actionUrl,
	}

	return p.publishEmail(ctx, "email.notification", message)
}

func (p *EmailProducer) SendEmailWarning(ctx context.Context, email, recipientName, content, actionUrl string) error {
	message := EmailMessage{
		Type:          "warning",
		Recipient:     email,
		RecipientName: recipientName,
		Content:       content,
		ActionUrl:     actionUrl,
	}

	return p.publishEmail(ctx, "email.warning", message)
}

func (p *EmailProducer) publishEmail(ctx context.Context, routingKey string, message EmailMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal email message: %w", err)
	}

	err = p.rabbitmq.Channel.PublishWithContext(
		ctx,
		"email_exchange", // exchange
		routingKey,       // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish email message: %w", err)
	}

	return nil
}
