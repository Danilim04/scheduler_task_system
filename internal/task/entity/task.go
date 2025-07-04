package entity

import (
	"context"
	"errors"
	"time"

	"github.com/robfig/cron"
)

type TaskRepositoryTemplateInterface interface {
	CreateTemplate(ctx context.Context, task *Task) error
}

type TaskRepositoryInterface interface {
	ExistsByID(ctx context.Context, id TaskID) (bool, error)
	Save(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id TaskID) (*Task, error)
	FindAll(ctx context.Context) ([]*Task, error)
	DeleteByID(ctx context.Context, id TaskID) error
	Update(ctx context.Context, task *Task) error
}

type TaskStatus string

const (
	TaskStatusActive   TaskStatus = "active"
	TaskStatusInactive TaskStatus = "inactive"
)

type TaskID string

type Task struct {
	TaskId      TaskID                 `json:"taskId" bson:"task_id"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Config      map[string]interface{} `json:"config" bson:"config"`
	Schedule    Schedule               `json:"schedule" bson:"schedule"`
	Status      TaskStatus             `json:"status" bson:"status"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

type Schedule struct {
	Expression string     `json:"expression" bson:"expression"`
	NextRun    *time.Time `json:"next_run,omitempty" bson:"next_run"`
	LastRun    *time.Time `json:"last_run,omitempty" bson:"last_run"`
}

type TaskExecution struct {
	TaskExecutionId string                 `json:"id" bson:"task_execution_id"`
	TaskID          TaskID                 `json:"task_id" bson:"task_id"`
	WorkerID        string                 `json:"worker_id" bson:"worker_id"`
	StartedAt       time.Time              `json:"started_at" bson:"started_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Success         bool                   `json:"success" bson:"success"`
	Error           string                 `json:"error,omitempty" bson:"error,omitempty"`
	Result          map[string]interface{} `json:"result,omitempty" bson:"result,omitempty"`
	RetryCount      int                    `json:"retry_count" bson:"retry_count"`
	Duration        time.Duration          `json:"duration" bson:"duration"`
}

func NewCreateTask(
	id TaskID,
	name string,
	description string,
	config map[string]interface{},
	expression string,
) (*Task, error) {
	task := &Task{
		TaskId:      id,
		Name:        name,
		Description: description,
		Config:      config,
		Schedule: Schedule{
			Expression: expression,
			NextRun:    nil,
			LastRun:    nil,
		},
		Status:    TaskStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := task.IsValid()

	if err != nil {
		return nil, err
	}

	return task, err
}

func (s *Task) IsValid() error {
	if s.Name == "" {
		return errors.New("invalid task name")
	}
	if s.Description == "" {
		return errors.New("invalid task description")
	}
	if len(s.Config) == 0 {
		return errors.New("invalid task config")
	}
	if s.Status != TaskStatusActive && s.Status != TaskStatusInactive {
		return errors.New("invalid task status")
	}
	if err := isValidCron(s.Schedule.Expression); err != nil {
		return err
	}
	return nil
}

func isValidCron(expr string) error {
	schedule, err := cron.ParseStandard(expr)

	if err != nil {
		return err
	}

	now := time.Now()
	next := schedule.Next(now)
	next2 := schedule.Next(next)

	if next2.Sub(next) >= 10*time.Minute {
		return nil
	}

	return err
}
