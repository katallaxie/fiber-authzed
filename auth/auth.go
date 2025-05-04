package auth

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3filter"
)

// OpenAPIAuthenticatorOpts are the OpenAPI authenticator options.
type OpenAPIAuthenticatorOpts struct {
	BasicAuthenticator openapi3filter.AuthenticationFunc
}

// Conigure the OpenAPI authenticator.
func (o *OpenAPIAuthenticatorOpts) Conigure(opts ...OpenAPIAuthenticatorOpt) {
	for _, opt := range opts {
		opt(o)
	}
}

// OpenAPIAuthenticatorOpt is a function that sets an option on the OpenAPI authenticator.
type OpenAPIAuthenticatorOpt func(*OpenAPIAuthenticatorOpts)

func OpenAPIAuthenticatorDefaultOpts() OpenAPIAuthenticatorOpts {
	return OpenAPIAuthenticatorOpts{}
}

// NoopBasicAuthenticator is a no-op basic authenticator.
func NoopBasicAuthenticator(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	return nil
}

// WithBasicAuthenticator sets the basic authenticator.
func WithBasicAuthenticator(auth openapi3filter.AuthenticationFunc) OpenAPIAuthenticatorOpt {
	return func(o *OpenAPIAuthenticatorOpts) {
		o.BasicAuthenticator = auth
	}
}

// NewAuthenticator returns a new authenticator.
func NewAuthenticator(opts ...OpenAPIAuthenticatorOpt) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		options := OpenAPIAuthenticatorDefaultOpts()
		options.Conigure(opts...)

		if input.SecuritySchemeName == "basic_auth" {
			return options.BasicAuthenticator(ctx, input)
		}

		return nil
	}
}
