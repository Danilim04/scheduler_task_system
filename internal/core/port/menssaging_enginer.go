package port

import (
	"context"
	"scheduler_task_system/internal/core/entity"
)

type MessagingEngine interface {
	Publish(ctx context.Context, task entity.TaskExecution) error
	Consumer(ctx context.Context) error
}
