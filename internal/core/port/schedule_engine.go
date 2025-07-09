package port

import (
	"context"
	"scheduler_task_system/internal/core/entity"

	"github.com/go-co-op/gocron/v2"
)

type ScheduleEngine interface {
	Register(ctx context.Context, taskName entity.TaskID, schedule entity.Schedule, payload []byte) (gocron.Job, error)
	Start() error
}
