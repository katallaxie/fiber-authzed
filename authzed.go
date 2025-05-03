package authzed

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

// Checker is the interface checking for authorization.
type Checker interface {
	// Allowed ...
	Allowed(context.Context) (bool, error)
}

// New creates a new middleware handler.
func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	// Return a new middleware handler
	return func(c fiber.Ctx) error {
		// Check if the next function is defined and returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Call the next handler in the chain
		return c.Next()
	}
}

// Helper function to set default values.
func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	if cfg.Checker == nil {
		cfg.Checker = ConfigDefault.Checker
	}

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = ConfigDefault.ErrorHandler
	}

	return cfg
}
