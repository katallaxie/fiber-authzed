package main

import (
	"context"
	"log"
	"os"

	"github.com/katallaxie/pkg/logx"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/spf13/cobra"
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
	Flags: &Flags{},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.Flags.Addr, "addr", ":8080", "addr")

	rootCmd.SilenceUsage = true
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, _ []string) error {
		return run(cmd.Context())
	},
}

func run(_ context.Context) error {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	_, err := logx.RedirectStdLog(logx.LogSink)
	if err != nil {
		return err
	}

	app := fiber.New()
	app.Use(requestid.New())

	err = app.Listen(cfg.Flags.Addr)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
