package util

import (
	"log"

	"github.com/jaysm12/multifinance-apps/pkg/rabbitmq"
)

func FlushQueue(queues []string, rabbitmqClient *rabbitmq.RabbitMqClient) error {
	for _, queue := range queues {
		_, err := rabbitmqClient.Channel.QueuePurge(queue, false)
		if err != nil {
			return err
		}
	}

	log.Println("Flush Queue Done")
	return nil
}
