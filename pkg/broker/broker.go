package broker

import "net/http"

type Broker interface {
	Publish(r *http.Request) error
	Consume() error
	Close()
}
