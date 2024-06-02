package broker

import "context"

type Broker interface {
	Publish(ctx context.Context, body string)
	Close()
}
