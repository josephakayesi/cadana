package main

import (
	// "encoding/json"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/gofiber/helmet/v2"
	route "github.com/josephakayesi/cadana/client/application/api/route"
	"github.com/josephakayesi/cadana/client/infra/config"

	slog "golang.org/x/exp/slog"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	c := config.NewConfig()

	app := fiber.New()

	app.Use(helmet.New())

	app.Use(requestid.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${ip} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: time.RFC3339Nano,
	}))

	timeout := time.Duration(time.Second*1) * time.Second

	route.Setup(timeout, app)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		log.Info("Shutting down gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Error("Server shutdown error:", err)
		}
	}()

	log.Info(fmt.Sprintf("server up and runing on port %d", c.PORT))

	err := app.Listen(fmt.Sprintf(":%d", c.PORT))
	if err != nil {
		panic(fmt.Sprintf("server was unable to start and listen on port %d", c.PORT))
	}
}
