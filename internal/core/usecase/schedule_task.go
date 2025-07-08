package usecase

import (
	"context"
	"errors"
	"scheduler_task_system/internal/core/entity"
	"scheduler_task_system/internal/core/port"
	"time"
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
	Task             entity.Task
	ProducerEnginner port.MessagingEngine
}

func (st *ScheduleTaskUc) Execute(ctx context.Context, input ScheduleInputDto) error {
	taskName := input.Task.TaskId
	schedule := input.Task.Schedule
	fn := func() error {
		taskexecutin := entity.TaskExecution{}
		execCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		uc := NewProducerTask(input.ProducerEnginner)
		input := ProducerTaskInput{
			TaskExecution: taskexecutin,
		}
		if err := uc.Execute(execCtx, input); err != nil {
			return err
		}
		return nil
	}
	if err := st.ScheduleEngine.Register(ctx, taskName, schedule, fn); err != nil {
		return errors.New("erro ao agendar a tarefa" + err.Error())
	}
	return nil
}

func (st *ScheduleTaskUc) ExecuteRemove(ctx context.Context, taskId entity.TaskID) error {
	if err := st.ScheduleEngine.Remove(ctx, taskId); err != nil {
		return errors.New("erro ao remover agendamento da tarefa" + err.Error())
	}
	return nil
}
