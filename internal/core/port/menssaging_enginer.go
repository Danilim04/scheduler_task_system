package port

import (
	"context"
	"scheduler_task_system/internal/core/entity"
)

type MessagingEngine interface {
	Publish(ctx context.Context, taskExecutin entity.TaskExecution, payload []byte) error
}
