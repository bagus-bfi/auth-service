package main

import (
	"log"

	"bravo-go-template/internal/config"
	"bravo-go-template/internal/event/subscriber"
	"bravo-go-template/internal/grpc"

	"github.com/bfi-finance/bfi-go-pkg/eventbus"
	"github.com/bfi-finance/bfi-go-pkg/logger"
	"github.com/bfi-finance/bfi-go-pkg/pgsql"
	"github.com/bfi-finance/bfi-go-pkg/tracer"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// LOAD APPLICATION CONFIG FROM ENVIRONMENT VARIABLES
	// -----------------------------------------------------------------------------------------------------------------
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("grpc: main failed to load and parse config: %s", err)
		return
	}
	// -----------------------------------------------------------------------------------------------------------------
	// STRUCTURED LOGGER
	// -----------------------------------------------------------------------------------------------------------------
	zlogger := logger.New(cfg.Logger).With().
		Str("app", cfg.AppInfo.Name).
		Str("env", cfg.AppInfo.Env).
		Str("rev", cfg.AppInfo.GitCommitHash).
		Logger()
	// -----------------------------------------------------------------------------------------------------------------
	// SET OPEN TELEMETRY GLOBAL TRACER
	// -----------------------------------------------------------------------------------------------------------------
	if err = tracer.SetTracer(cfg.Tracer, cfg.AppInfo); err != nil {
		zlogger.Error().Err(err).Msgf("grpc: main failed to setup open telemetry tracer: %s", err)
		return
	}
	// -----------------------------------------------------------------------------------------------------------------
	// INTERNAL EVENT BUS AND ITS BASE SUBSCRIBERS
	// -----------------------------------------------------------------------------------------------------------------
	eventBus := eventbus.NewInternal(zlogger)
	eventBus.RegisterSubscribers(
		subscriber.NewAppLogger(zlogger), // base application internal error logger.
	)
	// -----------------------------------------------------------------------------------------------------------------
	// INFRASTRUCTURE OBJECTS
	// -----------------------------------------------------------------------------------------------------------------
	// PGSQL
	sqlDB, sqlDBErr := pgsql.NewDB(cfg.PostgreSQL, zlogger)
	if sqlDBErr != nil {
		zlogger.Error().Err(sqlDBErr).Msgf("grpc: main failed to construct pgsql: %s", sqlDBErr)
		return
	}
	defer func() { _ = sqlDB.Close() }()
	// -----------------------------------------------------------------------------------------------------------------
	// GRPC SERVER SETUP
	// -----------------------------------------------------------------------------------------------------------------
	grpcServer, grpcServerErr := grpc.New(
		// ===== BASE DEPENDENCIES =====
		cfg, zlogger, eventBus,
		// ===== PER SERVICE DEPENDENCIES =====
		sqlDB,
	)
	if grpcServerErr != nil {
		zlogger.Error().Err(grpcServerErr).Msgf("grpc: main failed to construct server: %s", grpcServerErr)
		return
	}
	if err = grpcServer.Serve(); err != nil {
		zlogger.Error().Err(err).Msgf("grpc: main failed to serve: %s", err)
		return
	}
}
