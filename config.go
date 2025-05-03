package authzed

import (
	client "github.com/authzed/authzed-go/v1"
	"github.com/gofiber/fiber/v3"
)

type Client struct {
	client *client.Client
}

// NewClient creates a new Authzed client.
func NewClient(client *client.Client) *Client {
	return &Client{
		client: client,
	}
}

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c fiber.Ctx) bool
	// Checker is the implementation to check for authorization.
	Checker Checker
	// ErrorHandler is executed when an error is returned from fiber.Handler.
	//
	// Optional. Default: DefaultErrorHandler
	ErrorHandler fiber.ErrorHandler
}

// ConfigDefault is the default config.
var ConfigDefault = Config{
	Next: nil,
}
