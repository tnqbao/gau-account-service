package infra

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tnqbao/gau-account-service/shared/config"
)

type RabbitMQClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func InitRabbitMQClient(cfg *config.EnvConfig) *RabbitMQClient {
	rabbitUser := cfg.RabbitMQ.Username
	rabbitPassword := cfg.RabbitMQ.Password
	rabbitHost := cfg.RabbitMQ.Host
	rabbitPort := cfg.RabbitMQ.Port

	if rabbitUser == "" || rabbitPassword == "" || rabbitHost == "" || rabbitPort == "" {
		log.Fatal("One or more required RabbitMQ secrets are missing")
	}

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitUser, rabbitPassword, rabbitHost, rabbitPort,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatal("Failed to open a channel:", err)
	}

	log.Println("RabbitMQ connected at", rabbitHost)

	return &RabbitMQClient{
		Connection: conn,
		Channel:    ch,
	}
}

func (r *RabbitMQClient) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
	log.Println("RabbitMQ connection closed")
}

func (r *RabbitMQClient) DeclareQueue(queueName string, durable, autoDelete bool) error {
	_, err := r.Channel.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}
	log.Printf("Queue declared: %s", queueName)
	return nil
}

func (r *RabbitMQClient) DeclareExchange(exchangeName, exchangeType string, durable bool) error {
	err := r.Channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type (direct, fanout, topic, headers)
		durable,      // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", exchangeName, err)
	}
	log.Printf("Exchange declared: %s (type: %s)", exchangeName, exchangeType)
	return nil
}

func (r *RabbitMQClient) BindQueue(queueName, exchangeName, routingKey string) error {
	err := r.Channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueName, exchangeName, err)
	}
	log.Printf("Queue %s bound to exchange %s with routing key %s", queueName, exchangeName, routingKey)
	return nil
}
