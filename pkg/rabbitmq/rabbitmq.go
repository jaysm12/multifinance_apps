package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQClient(config interface{}) (*RabbitMqClient, error) {
	conn, err := amqp.Dial(config.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return &RabbitMqClient{
		Connection: conn,
		Channel:    ch,
	}, nil
}
