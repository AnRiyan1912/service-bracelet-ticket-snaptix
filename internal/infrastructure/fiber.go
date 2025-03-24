package infrastructure

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	braceletticket "bracelet-ticket-system-be/internal/bracelet-ticket"
	"bracelet-ticket-system-be/internal/middleware"
	"bracelet-ticket-system-be/pkg/xlogger"
)

func Run() {
	logger := xlogger.Logger
	app := fiber.New(fiber.Config{
		ProxyHeader:             cfg.ProxyHeader,
		DisableStartupMessage:   true,
		ErrorHandler:            middleware.DefaultErrorHandler,
		EnableTrustedProxyCheck: true,
	})

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
		Fields: cfg.LogFields,
	}))

	runWebsocket(app, mysqlEventRepository, redisEventRepository, redisClient)

	api := app.Group("/api/v1.0")
	braceletticket.NewHttpHandler(api.Group("/bracelet-tickets"), braceletTicketService)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		logger.Info().Msg("Shutting down gracefully...")

		// Stop the Fiber
		if err := app.Shutdown(); err != nil {
			logger.Error().Err(err).Msg("Error shutting down server")
		}

		logger.Info().Msg("Server stopped")
	}()

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info().Msgf("Server is running on address: %s", addr)
	if err := app.Listen(addr); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}

}
