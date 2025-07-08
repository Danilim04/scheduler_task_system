package shedulerenginner

import (
	"context"
	"scheduler_task_system/internal/core/entity"
	"scheduler_task_system/internal/core/port"
	"scheduler_task_system/internal/core/usecase"
	"time"
)

var ProducerEnginner port.MessagingEngine

func ExecutarTask(taskname entity.TaskID, expression string, payload []byte) error {

	taskExecutin := entity.TaskExecution{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uc := usecase.NewProducerTask(ProducerEnginner)

	input := usecase.ProducerTaskInput{
		TaskExecution: taskExecutin,
		Payload:       payload,
	}

	if err := uc.Execute(ctx, input); err != nil {
		return err
	}
	return nil
}
