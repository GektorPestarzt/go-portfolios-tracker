package main

import (
	"fmt"
	"go-portfolios-tracker/internal/config"
	"go-portfolios-tracker/internal/logging/slog"
	"go-portfolios-tracker/pkg/broker"
	"go-portfolios-tracker/pkg/broker/rabbitmq"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.server(w, r)
		return
	}

	err := h.broker.Publish(r)
	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(201)
	}
}

func (h *Handler) server(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest(r.Method, r.URL.String(), r.Body)
	for key, value := range r.Header {
		req.Header.Set(key, strings.Join(value, " "))
	}

	req.URL.Scheme = "http"
	req.URL.Host = "localhost:1234"

	client := http.Client{}
	resp, _ := client.Do(req)

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.Write(respBody)
}

func scanner(broker broker.Broker) {
	handler := NewHandler(broker)

	http.HandleFunc("/api/accounts/tinkoff", handler.update)
	http.HandleFunc("/auth", handler.server)
	http.ListenAndServe(":8080", nil)
}
