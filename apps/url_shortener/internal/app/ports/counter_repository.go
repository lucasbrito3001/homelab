package ports

import "context"

type CounterRepository interface {
	Increment(ctx context.Context) (int64, error)
}
