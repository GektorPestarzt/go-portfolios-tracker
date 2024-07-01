package rabbitmq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-portfolios-tracker/internal/logging"
	"io/ioutil"
	"net/http"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	RequestQueue amqp.Queue

	logger logging.Logger
}

type HTTPRequestData struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type HTTPResponseData struct {
	StatusCode int               `json:"statuscode"`
	Headers    map[string]string `json:headers`
	Body       string            `json:body`
}

func NewRabbitMQ(logger logging.Logger) *RabbitMQ {
	uri := fmt.Sprintf("amqp://%s:%s@localhost:5672", "lab", "password")

	conn, err := amqp.Dial(uri)
	if err != nil {
		logger.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("failed to open channel. Error: %s", err)
	}

	request, err := ch.QueueDeclare(
		"request",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatalf("failed to declare a queue. Error: %s", err)
	}

	return &RabbitMQ{
		Connection:   conn,
		Channel:      ch,
		RequestQueue: request,
		logger:       logger,
	}
}

func (r *RabbitMQ) Publish(req *http.Request) error {
	data, _ := extractRequestData(req)
	jsonData, _ := json.Marshal(data)

	err := r.Channel.Publish(
		"",
		r.RequestQueue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonData,
		})
	if err != nil {
		r.logger.Error("Failed to publish a message", err)
		return err
	}
	r.logger.Infof(" [x] Sent %s", jsonData)

	return nil
}

func (r *RabbitMQ) Consume() error {
	msgs, err := r.Channel.Consume(
		r.RequestQueue.Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		return err
	}

	var forever chan struct{}
	// 	ctx := context.WithoutCancel(context.Background())

	go func() {
		for d := range msgs {
			r.logger.Infof("Received a message")

			var request HTTPRequestData
			json.Unmarshal(d.Body, &request)

			req, _ := http.NewRequest(request.Method, request.URL, bytes.NewBuffer([]byte(request.Body)))
			for key, value := range request.Headers {
				req.Header.Set(key, value)
			}

			req.URL.Scheme = "http"
			req.URL.Host = "localhost:1234"

			client := http.Client{}
			client.Do(req)

			r.logger.Info("Done")
		}
	}()

	r.logger.Info(" [*] Waiting for messages.")
	<-forever

	return nil
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Connection.Close()
}

func extractRequestData(r *http.Request) (*HTTPRequestData, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, " ")
	}

	return &HTTPRequestData{
		Method:  r.Method,
		URL:     r.URL.String(),
		Headers: headers,
		Body:    string(bodyBytes),
	}, nil
}

func extractResponseData(r *http.Response) (*HTTPResponseData, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, " ")
	}

	return &HTTPResponseData{
		StatusCode: r.StatusCode,
		Headers:    headers,
		Body:       string(bodyBytes),
	}, nil
}
