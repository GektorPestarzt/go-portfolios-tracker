package main

import (
	portfoliorepo "go-portfolios-tracker/internal/account/repository/sqlite"
	portfoliousecase "go-portfolios-tracker/internal/account/usecase/tinkoff"
	"go-portfolios-tracker/internal/config"
	"go-portfolios-tracker/internal/logging/slog"
	"go-portfolios-tracker/pkg/broker/rabbitmq"
	"go-portfolios-tracker/pkg/repository/sqlite"
	"os"
)

func main() {
	os.Setenv("CONFIG_PATH", "./config.yaml")
	cfg := config.MustLoad()

	logger := slog.NewLogger(cfg.Env)

	logger.Info("create rabitmq")
	broker := rabbitmq.NewRabbitMQ(logger)

	logger.Info("create database")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Fatal("failed to init repository", err)
	}
	portfolioRepo := portfoliorepo.NewPortfolioRepository(storage)

	portfolioUseCase := portfoliousecase.NewPortfolioUseCase(logger, portfolioRepo)
	broker.Consume(portfolioUseCase)
}
