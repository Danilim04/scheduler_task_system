package usecase

import (
	"context"
	"errors"
	"scheduler_task_system/internal/core/entity"
	"scheduler_task_system/internal/core/port"

	"github.com/go-co-op/gocron/v2"
)

type ScheduleTaskUc struct {
	ScheduleEngine port.ScheduleEngine
}

func NewScheduleTask(task entity.Task, f func(), se port.ScheduleEngine) *ScheduleTaskUc {
	return &ScheduleTaskUc{
		ScheduleEngine: se,
	}
}

type ScheduleInputDto struct {
	Task entity.Task
}

func (st *ScheduleTaskUc) Execute(ctx context.Context, input ScheduleInputDto) (gocron.Job, error) {

	taskName := input.Task.TaskId
	payload := input.Task.Payload
	schedule := input.Task.Schedule
	j, err := st.ScheduleEngine.Register(ctx, taskName, schedule, payload)
	if err != nil {
		return nil, errors.New("erro ao agendar a tarefa" + err.Error())
	}
	return j, nil
}
