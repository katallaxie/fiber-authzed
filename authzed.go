package authzed

import (
	"context"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	client "github.com/authzed/authzed-go/v1"
	"github.com/gofiber/fiber/v3"
)

// Checker is the interface checking for authorization.
type Checker interface {
	// Allowed ...
	Allowed(context.Context) (bool, error)
}

// Client is the Authzed client.
type Client struct {
	client *client.Client
}

// Allowed checks if the user is allowed to perform the action.
func (c *Client) Allowed(ctx context.Context, action string) (bool, error) {
	subj := &v1.SubjectReference{
		Object: &v1.ObjectReference{
			ObjectType: "blog/user",
			ObjectId:   "emilia",
		},
	}

	res := &v1.ObjectReference{
		ObjectType: "blog/post",
		ObjectId:   "1",
	}

	req := &v1.CheckPermissionRequest{
		Resource:   res,
		Permission: action,
		Subject:    subj,
	}

	resp, err := c.client.CheckPermission(ctx, req)
	if err != nil {
		return false, err
	}

	return resp.GetPermissionship() == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION, nil
}

// NewClient creates a new Authzed client.
func NewClient(client *client.Client) *Client {
	return &Client{
		client: client,
	}
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
