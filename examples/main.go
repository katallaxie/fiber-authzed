package main

import (
	"context"
	"log"
	"os"

	"github.com/katallaxie/fiber-authzed/auth"
	"github.com/katallaxie/fiber-authzed/examples/api"
	"github.com/katallaxie/fiber-authzed/oas"

	client "github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	authzed "github.com/katallaxie/fiber-authzed"
	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/pkg/logx"
	"github.com/katallaxie/pkg/server"
	middleware "github.com/oapi-codegen/fiber-middleware"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config ...
type Config struct {
	Flags *Flags
}

// Flags ...
type Flags struct {
	Addr string
}

var cfg = &Config{
	Flags: &Flags{
		Addr: ":8080",
	},
}

// WebSrv is the server that implements the Noop interface.
type WebSrv struct {
	cfg *Config
}

// NewWebSrv returns a new instance of NoopSrv.
func NewWebSrv(cfg *Config) *WebSrv {
	return &WebSrv{cfg}
}

var _ api.StrictServerInterface = (*apiHandlers)(nil)

type apiHandlers struct{}

// NewAPIHandlers returns a new instance of APIHandlers.
func NewAPIHandlers() *apiHandlers {
	return &apiHandlers{}
}

// Echo is a simple echo handler.
func (h *apiHandlers) Echo(ctx context.Context, request api.EchoRequestObject) (api.EchoResponseObject, error) {
	return api.Echo200JSONResponse{Echo: cast.Ptr("hello world")}, nil
}

// Start starts the server.
func (s *WebSrv) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		swagger, err := api.GetSwagger()
		if err != nil {
			return err
		}
		swagger.Servers = nil

		app := fiber.New()
		app.Use(requestid.New())
		app.Use(logger.New())

		c, err := client.NewClient(
			"localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpcutil.WithInsecureBearerToken("example"),
		)
		if err != nil {
			return err
		}
		check := authzed.NewChecker(c)

		validatorOptions := &middleware.Options{}
		validatorOptions.Options.AuthenticationFunc = auth.NewAuthenticator(auth.WithBasicAuthenticator(auth.NoopBasicAuthenticator))
		validatorOptions.Options.AuthenticationFunc = oas.Authenticate(
			oas.OasAuthenticate(
				oas.WithChecker(check),
			),
		)

		app.Use(middleware.OapiRequestValidatorWithOptions(swagger, validatorOptions))

		handlers := NewAPIHandlers()
		handler := api.NewStrictHandler(handlers, nil)
		api.RegisterHandlers(app, handler)

		err = app.Listen(s.cfg.Flags.Addr)
		if err != nil {
			return err
		}

		return nil
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.Flags.Addr, "addr", cfg.Flags.Addr, "addr")

	rootCmd.SilenceUsage = true
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, _ []string) error {
		return run(cmd.Context())
	},
}

func run(ctx context.Context) error {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	_, err := logx.RedirectStdLog(logx.LogSink)
	if err != nil {
		return err
	}

	srv := NewWebSrv(cfg)

	s, _ := server.WithContext(ctx)
	s.Listen(srv, false)

	return s.Wait()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
