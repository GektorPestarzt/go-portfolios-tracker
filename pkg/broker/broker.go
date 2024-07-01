package broker

import "net/http"

type Broker interface {
	Publish(r *http.Request) []byte
	Consume() error
	Close()
}
