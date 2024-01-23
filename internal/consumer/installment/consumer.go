package installment

import (
	"encoding/json"
	"log"

	"github.com/jaysm12/multifinance-apps/internal/service/installment"
	"github.com/jaysm12/multifinance-apps/pkg/rabbitmq"
)

type InstallmentConsumerMethod interface {
	CreateInstallmentConsumer(errChan chan error) error
	PayInstallmentConsumer(errChan chan error) error
}

const (
	queueCreateInstallment = "create_installment"
	queuePayInstallment    = "pay_installment"
)

type InstallmentConsumer struct {
	InstallmentService installment.InstallmentServiceMethod
	rmqClient          *rabbitmq.RabbitMqClient
}

func NewInstallmentConsumer(installmentService installment.InstallmentServiceMethod, rmqClient *rabbitmq.RabbitMqClient) InstallmentConsumerMethod {
	return &InstallmentConsumer{
		InstallmentService: installmentService,
		rmqClient:          rmqClient,
	}
}

func (i *InstallmentConsumer) CreateInstallmentConsumer(errChan chan error) error {
	// declare queue
	q, err := i.rmqClient.Channel.QueueDeclare(
		queueCreateInstallment, // name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)

	if err != nil {
		log.Fatal(err, "Failed to declare a queue")
		errChan <- err
	}
	msgs, err := i.rmqClient.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err, "Failed to register a consumer create_installment")
		errChan <- err
	}

	for msg := range msgs {
		var body CreateInstallmentPayload
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			log.Printf("Failed to unmarshal JSON message: %v", err)
			errChan <- err
			continue
		}

		// Process the JSON message
		err = i.InstallmentService.CreateInstallment(installment.CreateInstallmentRequest{
			UserID:         body.UserID,
			CreditOptionID: body.CreditOptionID,
			OtrAmount:      body.OtrAmount,
			AssetName:      body.AssetName,
		})

		if err != nil {
			log.Printf("Failed to process message: %v", err)
			errChan <- err
		}
	}
	return nil
}

func (i *InstallmentConsumer) PayInstallmentConsumer(errChan chan error) error {
	// declare queue
	q, err := i.rmqClient.Channel.QueueDeclare(
		queuePayInstallment, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)

	if err != nil {
		log.Fatal(err, "Failed to declare a queue")
		errChan <- err
	}
	msgs, err := i.rmqClient.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err, "Failed to register a consumer create_installment")
		errChan <- err
	}

	for msg := range msgs {
		var body PayInstallmentPayload
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			log.Printf("Failed to unmarshal JSON message: %v", err)
			errChan <- err
			continue
		}

		// Process the JSON message
		err = i.InstallmentService.PayInstallment(installment.PayInstallmentRequest{
			UserID:     body.UserID,
			PaidAmount: body.PaidAmount,
			ContractID: body.ContractID,
		})

		if err != nil {
			log.Printf("Failed to process message: %v", err)
			errChan <- err
		}
	}
	return nil
}
