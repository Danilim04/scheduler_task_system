package shedulerenginner

import (
	"context"
	"scheduler_task_system/internal/core/entity"

	"github.com/go-co-op/gocron/v2"
)

type GoCron struct {
	GoCron gocron.Scheduler
}

func NewShedulerEngginer() (*GoCron, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &GoCron{
		GoCron: s,
	}, nil
}

func (g *GoCron) Register(ctx context.Context, taskName entity.TaskID, schedule entity.Schedule, payload []byte) (gocron.Job, error) {
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

func (g *GoCron) Start() error {
	if err := g.Start(); err != nil {
		return nil
	}
	return nil
}
