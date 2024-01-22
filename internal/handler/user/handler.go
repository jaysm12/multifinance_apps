package user

import (
	"github.com/jaysm12/multifinance-apps/internal/service/user"
)

type UserHandler struct {
	service      user.UserServiceMethod
	timeoutInSec int
}

// Option set options for http handler config
type Option func(*UserHandler)

const (
	defaultTimeout = 5
)

// NewUserHandler is func to create http auth handler
func NewUserHandler(service user.UserServiceMethod, options ...Option) *UserHandler {
	handler := &UserHandler{
		service:      service,
		timeoutInSec: defaultTimeout,
	}

	// Apply options
	for _, opt := range options {
		opt(handler)
	}

	return handler
}

// WithTimeoutOptions is func to set timeout config into handler
func WithTimeoutOptions(timeoutinsec int) Option {
	return Option(
		func(h *UserHandler) {
			if timeoutinsec <= 0 {
				timeoutinsec = defaultTimeout
			}
			h.timeoutInSec = timeoutinsec
		})
}
