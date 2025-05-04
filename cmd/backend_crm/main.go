package main

import (
	"backend_crm/internal/config"
	httpController "backend_crm/internal/controller/http/fasthttp"
	"backend_crm/internal/controller/http/fasthttp/app"
	"backend_crm/internal/controller/http/fasthttp/authorization"
	"backend_crm/internal/controller/http/fasthttp/orders"
	ordersRepo "backend_crm/internal/repository/orders/postgre"
	usersRepo "backend_crm/internal/repository/users/postgre"
	"backend_crm/internal/usecase/users/std"
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

func main() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Initialize database connection
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping database")
	}

	// Initialize repositories
	usersRepo := usersRepo.NewRepository(db)
	ordersRepo := ordersRepo.NewRepository(db)

	// Initialize usecases
	usersUsecase := std.NewUsecase(
		usersRepo,
		[]byte(cfg.JWT.AccessSecret),
		[]byte(cfg.JWT.RefreshSecret),
		cfg.GetAccessTTL(),
		cfg.GetRefreshTTL(),
	)

	// Initialize controllers
	authController := authorization.NewController(usersUsecase, logger.With().Str("component", "authorization").Logger())
	ordersController := orders.NewController(ordersRepo, logger.With().Str("component", "orders").Logger())
	appController := app.NewController(cfg.HTML.Files.Index, logger.With().Str("component", "app").Logger())

	// Initialize main controller
	controller := httpController.NewController(
		*authController,
		*ordersController,
		*appController,
	)

	// Create server
	server := &fasthttp.Server{
		Handler:            controller.Handlers(context.Background()),
		ReadTimeout:        cfg.GetReadTimeout(),
		WriteTimeout:       cfg.GetWriteTimeout(),
		MaxRequestBodySize: 10 * 1024 * 1024, // 10MB
	}

	// Create error channel
	errChan := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		logger.Info().Str("addr", cfg.GetServerAddr()).Msg("starting server")
		if err := server.ListenAndServeTLS(cfg.GetServerAddr(), cfg.TLS.CertFilePath, cfg.TLS.CertKeyPath); err != nil {
			errChan <- err
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Error().Err(err).Msg("server error")
	case sig := <-quit:
		logger.Info().Str("signal", sig.String()).Msg("received signal")
	}

	// Graceful shutdown
	logger.Info().Msg("shutting down server")
	if err := server.Shutdown(); err != nil {
		logger.Error().Err(err).Msg("error during server shutdown")
	}
}
