package rabbitmq

import (
	"context"
	"fmt"
	"go-portfolios-tracker/internal/account"
	"go-portfolios-tracker/internal/logging"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue

	logger logging.Logger
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

	q, err := ch.QueueDeclare(
		"update",
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
		Connection: conn,
		Channel:    ch,
		Queue:      q,
		logger:     logger,
	}
}

func (r *RabbitMQ) Publish(ctx context.Context, body string) {
	err := r.Channel.PublishWithContext(ctx,
		"",
		r.Queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})

	if err != nil {
		r.logger.Error("Failed to publish a message", err)
	}
	r.logger.Infof(" [x] Sent %s", body)
}

func (r *RabbitMQ) Consume(useCase account.UseCase) error {
	msgs, err := r.Channel.Consume(
		r.Queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		return err
	}

	var forever chan struct{}
	ctx := context.WithoutCancel(context.Background())

	go func() {
		for d := range msgs {
			r.logger.Infof("Received a message: %s", d.Body)

			id, _ := strconv.Atoi(string(d.Body))
			useCase.Update(ctx, int64(id))

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
