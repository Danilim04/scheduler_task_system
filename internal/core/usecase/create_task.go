package usecase

import (
	"context"
	"errors"
	"scheduler_task_system/internal/core/entity"
	"scheduler_task_system/internal/core/port"

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
	TaskRepository         port.TaskRepositoryInterface
	TaskRepositoryTemplate port.TaskRepositoryTemplateInterface
}

func NewCreateTaskUseCase(taskRepository port.TaskRepositoryInterface, taskRepositoryTemplate port.TaskRepositoryTemplateInterface) *CreateTaskUseCase {
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
		return nil, errors.New("falha ao criar a task" + err.Error())
	}

	if err := uc.TaskRepositoryTemplate.Generate(ctx, task); err != nil {
		return nil, errors.New("falha ao criar o código de template da task" + err.Error())
	}

	if err := uc.TaskRepository.Save(ctx, task); err != nil {
		return nil, errors.New("falha ao criar task no banco de dados" + err.Error())
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
