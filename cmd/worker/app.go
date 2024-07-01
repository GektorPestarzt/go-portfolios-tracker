package main

import (
	"fmt"
	"go-portfolios-tracker/internal/config"
	"go-portfolios-tracker/internal/logging/slog"
	"go-portfolios-tracker/pkg/broker/rabbitmq"
	"os"
)

func main() {
	os.Setenv("CONFIG_PATH", "./config.yaml")
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger := slog.NewLogger(cfg.Env)

	broker := rabbitmq.NewRabbitMQ(logger)
	defer broker.Close()
	broker.Consume()
}
