package main

import (
	"encoding/json"
	"fmt"
	"go-portfolios-tracker/internal/config"
	"go-portfolios-tracker/internal/logging/slog"
	"go-portfolios-tracker/pkg/broker"
	"go-portfolios-tracker/pkg/broker/rabbitmq"
	"net/http"
	"os"
)

func main() {
	os.Setenv("CONFIG_PATH", "./config.yaml")
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger := slog.NewLogger(cfg.Env)

	broker := rabbitmq.NewRabbitMQ(logger)
	defer broker.Close()
	scanner(broker)
}

type Handler struct {
	broker broker.Broker
}

func NewHandler(broker broker.Broker) *Handler {
	return &Handler{broker: broker}
}

func (h *Handler) server(w http.ResponseWriter, r *http.Request) {
	respBytes := h.broker.Publish(r)
	// var resp http.Response
	// json.Unmarshal(respBytes, &resp)
	var resp rabbitmq.HTTPResponseData
	json.Unmarshal(respBytes, &resp)

	for name, values := range resp.Headers {
		w.Header().Add(name, values)
	}
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(resp.Body))
}

func scanner(broker broker.Broker) {
	handler := NewHandler(broker)

	http.HandleFunc("/", handler.server)
	http.ListenAndServe(":8080", nil)
}
