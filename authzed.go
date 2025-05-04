package authzed

import (
	"context"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	client "github.com/authzed/authzed-go/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/pkg/conv"
)

// Checker is the interface checking for authorization.
type Checker interface {
	// Allowed ...
	Allowed(ctx context.Context, resource *ObjectReference, subject *SubjectReference, permission *Permission) (bool, error)
}

var _ Checker = (*Client)(nil)

// Client is the Authzed client.
type Client struct {
	client *client.Client
}

// SubjectReference is used to refer to a specific subject in the system.
type SubjectReference struct {
	Object *ObjectReference
}

// ObjectReference is used to refer to a specific object in the system.
type ObjectReference struct {
	ObjectType string
	ObjectId   string
}

// Permission is the permission to check.
type Permission string

// Allowed checks if the user is allowed to perform the action.
func (c *Client) Allowed(ctx context.Context, resource *ObjectReference, subject *SubjectReference, permission *Permission) (bool, error) {
	subj := &v1.SubjectReference{
		Object: &v1.ObjectReference{
			ObjectType: subject.Object.ObjectType,
			ObjectId:   subject.Object.ObjectId,
		},
	}

	res := &v1.ObjectReference{
		ObjectType: resource.ObjectType,
		ObjectId:   resource.ObjectId,
	}

	req := &v1.CheckPermissionRequest{
		Resource:   res,
		Permission: conv.String(cast.Value(permission)),
		Subject:    subj,
	}

	resp, err := c.client.CheckPermission(ctx, req)
	if err != nil {
		return false, err
	}

	return resp.GetPermissionship() == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION, nil
}

// NewChecker creates a new Authzed client.
func NewChecker(client *client.Client) *Client {
	return &Client{
		client: client,
	}
}

// New creates a new middleware handler.
func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	// Return a new middleware handler
	return func(c *fiber.Ctx) error {
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
