package installment

import (
	"github.com/jaysm12/multifinance-apps/internal/service/installment"
	"github.com/jaysm12/multifinance-apps/pkg/rabbitmq"
)

type InstallmentHandler struct {
	service        installment.InstallmentServiceMethod
	timeoutInSec   int
	rabbitmqClient *rabbitmq.RabbitMqClient
}

// Option set options for http handler config
type Option func(*InstallmentHandler)

const (
	defaultTimeout = 5
)

// NewInstallmentHandler is func to create http auth handler
func NewInstallmentHandler(service installment.InstallmentServiceMethod, rabbitmqClient *rabbitmq.RabbitMqClient, options ...Option) *InstallmentHandler {
	handler := &InstallmentHandler{
		service:        service,
		timeoutInSec:   defaultTimeout,
		rabbitmqClient: rabbitmqClient,
	}

	// Apply options
	for _, opt := range options {
		opt(handler)
	}

	return handler
}

func WithTimeoutOptions(timeoutinsec int) Option {
	return Option(
		func(h *InstallmentHandler) {
			if timeoutinsec <= 0 {
				timeoutinsec = defaultTimeout
			}
			h.timeoutInSec = timeoutinsec
		})
}
