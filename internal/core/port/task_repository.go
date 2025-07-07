package port

import (
	"context"
	"scheduler_task_system/internal/core/entity"
)

type TaskRepositoryInterface interface {
	ExistsByID(ctx context.Context, id entity.TaskID) (bool, error)
	Save(ctx context.Context, task *entity.Task) error
	FindByID(ctx context.Context, id entity.TaskID) (*entity.Task, error)
	FindAll(ctx context.Context) ([]*entity.Task, error)
	DeleteByID(ctx context.Context, id entity.TaskID) error
	Update(ctx context.Context, task *entity.Task) error
}
