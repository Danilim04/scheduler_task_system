package shedulerenginner

import (
	"context"
	"scheduler_task_system/internal/core/entity"

	"github.com/go-co-op/gocron/v2"
)

type ScheduleEngine struct {
	GoCron gocron.Scheduler
}

func NewShedulerEngginer() (*ScheduleEngine, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &ScheduleEngine{
		GoCron: s,
	}, nil
}

func (g *ScheduleEngine) Register(ctx context.Context, taskName entity.TaskID, schedule entity.Schedule, payload []byte) (gocron.Job, error) {
	task := gocron.NewTask(
		ExecutarTask,
		taskName,
		payload,
	)
	j, err := g.GoCron.NewJob(gocron.CronJob(schedule.Expression, true), task)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (g *ScheduleEngine) Start() error {
	if err := g.Start(); err != nil {
		return nil
	}
	return nil
}
