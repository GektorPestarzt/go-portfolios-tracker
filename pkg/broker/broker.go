package broker

import (
	"context"
	"go-portfolios-tracker/internal/account"
)

type Broker interface {
	Publish(ctx context.Context, body string)
	Consume(useCase account.UseCase) error
	Close()
}
