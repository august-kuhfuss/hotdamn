package tasks

import "context"

type Task interface {
	Start(ctx context.Context) error
}
