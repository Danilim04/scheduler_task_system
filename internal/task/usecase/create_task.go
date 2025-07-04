package usecase

import (
	"context"
	"scheduler_task_system/internal/task/entity"
	"time"
)

type CreateTaskOutputDto struct {
	TaskId      entity.TaskID          `json:"taskId"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Schedule    entity.Schedule        `json:"schedule"`
	Status      entity.TaskStatus      `json:"status"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

type CreateTaskInputDto struct {
	TaskId      string                 `json:"taskId"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Expression  string                 `json:"expression"`
}

type CreateTaskUseCase struct {
	TaskRepository         entity.TaskRepositoryInterface
	TaskRepositoryTemplate entity.TaskRepositoryTemplateInterface
}

func NewCreateTaskUseCase(taskRepository entity.TaskRepositoryInterface, taskRepositoryTemplate entity.TaskRepositoryTemplateInterface) *CreateTaskUseCase {
	return &CreateTaskUseCase{TaskRepository: taskRepository, TaskRepositoryTemplate: taskRepositoryTemplate}
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, input CreateTaskInputDto) (*CreateTaskOutputDto, error) {

	task, err := entity.NewCreateTask(
		entity.TaskID(input.TaskId),
		input.Name,
		input.Description,
		input.Config,
		input.Expression,
	)

	if err != nil {
		return nil, err
	}

	err = uc.TaskRepositoryTemplate.CreateTemplate(ctx, task)

	if err != nil {
		return nil, err
	}

	return &CreateTaskOutputDto{
		TaskId:      task.TaskId,
		Name:        task.Name,
		Description: task.Description,
		Config:      task.Config,
		Schedule:    task.Schedule,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}
