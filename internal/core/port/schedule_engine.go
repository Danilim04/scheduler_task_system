package port

import (
	"context"
	"scheduler_task_system/internal/core/entity"
)

type ScheduleEngine interface {
	Register(ctx context.Context, taskName entity.TaskID, schedule entity.Schedule, fn func() error) error
	Remove(ctx context.Context, taskName entity.TaskID) error
	Start(ctx context.Context) error
}
