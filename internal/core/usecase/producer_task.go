package usecase

import (
	"context"
	"errors"
	"scheduler_task_system/internal/core/entity"
	"scheduler_task_system/internal/core/port"
)

type ProducerTask struct {
	ProducerEnginner port.MessagingEngine
}

func NewProducerTask(pe port.MessagingEngine) *ProducerTask {
	return &ProducerTask{
		ProducerEnginner: pe,
	}
}

type ProducerTaskInput struct {
	TaskExecution entity.TaskExecution
	Payload       []byte
}

func (pe *ProducerTask) Execute(ctx context.Context, input ProducerTaskInput) error {
	if err := pe.ProducerEnginner.Publish(ctx, input.TaskExecution, input.Payload); err != nil {
		return errors.New("erro ao produzir tarefa para a mensageria" + err.Error())
	}
	return nil
}
