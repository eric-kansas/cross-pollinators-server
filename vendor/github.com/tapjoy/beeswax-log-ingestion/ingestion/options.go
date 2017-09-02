package ingestion

import (
	"errors"
	"time"

	newrelic "github.com/newrelic/go-agent"
)

// Option Errors
var (
	ErrInvalidTimeout   = errors.New("Invalid Timeout")
	ErrInvalidRateLimit = errors.New("Invalid Rate Limit")
	ErrNoNewRelicApp    = errors.New("No new relic app provided")
)

// Option ...
type Option func(*server) error

// ReadTimeout ...
func ReadTimeout(sec int) Option {
	return func(s *server) error {
		if sec < 0 {
			return ErrInvalidTimeout
		}
		s.httpServer.ReadTimeout = time.Duration(sec) * time.Second
		return nil
	}
}

// WriteTimeout ...
func WriteTimeout(sec int) Option {
	return func(s *server) error {
		if sec < 0 {
			return ErrInvalidTimeout
		}
		s.httpServer.WriteTimeout = time.Duration(sec) * time.Second
		return nil
	}
}

// IdleTimeout ...
func IdleTimeout(sec int) Option {
	return func(s *server) error {
		if sec < 0 {
			return ErrInvalidTimeout
		}
		s.httpServer.IdleTimeout = time.Duration(sec) * time.Second
		return nil
	}
}

// NewRelic ...
func NewRelic(app *newrelic.Application) Option {
	return func(s *server) error {
		if app == nil {
			return ErrNoNewRelicApp
		}
		s.newrelicApp = *app
		return nil
	}
}

// RateLimit ...
func RateLimit(nano int64) Option {
	return func(s *server) error {
		if nano < 0 {
			return ErrInvalidRateLimit
		}
		s.rateLimit = nano
		return nil
	}
}

// ClientTimeout ...
func ClientTimeout(sec int) Option {
	return func(s *server) error {
		if sec < 0 {
			return ErrInvalidTimeout
		}
		httpClient.Timeout = time.Duration(sec) * time.Second
		return nil
	}
}
