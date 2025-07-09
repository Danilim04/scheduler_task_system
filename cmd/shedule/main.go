package main

import (
	"context"
	"log"
	"scheduler_task_system/internal/core/usecase"
	"scheduler_task_system/internal/infra/mongodb"
	shedulerenginner "scheduler_task_system/internal/infra/sheduler_enginner"
	"time"
)

func main() {

	ctx := context.Background()

	client, err := mongodb.ConnectMongodb()
	if err != nil {
		log.Fatal(err)
	}
	repository := mongodb.NewTaskRepositoryMongo(client)
	tasks, err := repository.FindAll(ctx)
	if err != nil {
		log.Fatal(err)
	}
	scheduleEngine, err := shedulerenginner.NewShedulerEngginer()
	if err != nil {
		log.Fatal(err)
	}
	scheduletask := usecase.NewScheduleTask(scheduleEngine)
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		input := usecase.ScheduleInputDto{
			Task: *task,
		}
		job, err := scheduletask.Execute(ctx, input)
		if err != nil {
			log.Fatal(err)
		}
		idJob := job.ID()
		task.Schedule.IdJob = &idJob
		if next, _ := job.NextRun(); next != (time.Time{}) {
			task.Schedule.NextRun = &next
		}
		task.UpdatedAt = time.Now()

		if err := task.IsValid(); err != nil {
			log.Fatal(err)
		}

		if err := repository.Update(ctx, task); err != nil {
			log.Printf("falha ao salvar task %s: %v", task.TaskId, err)
		}
	}
	scheduleEngine.Start()
}
