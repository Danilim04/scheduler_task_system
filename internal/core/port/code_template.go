package port

import (
	"context"
	"scheduler_task_system/internal/core/entity"
)

type TaskRepositoryTemplateInterface interface {
	Generate(ctx context.Context, task *entity.Task) error
}
