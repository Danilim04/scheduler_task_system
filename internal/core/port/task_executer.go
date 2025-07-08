package port

import "context"

type TaskExecuter interface {
	Execute(ctx context.Context, payload []byte) error
}
