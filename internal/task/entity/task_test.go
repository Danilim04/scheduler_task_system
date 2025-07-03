package entity_test

import (
	"testing"
	"time"

	"scheduler_task_system/internal/task/entity"

	"github.com/stretchr/testify/assert"
)

func TestGivenAnEmptyName_WhenCreateANewTask_ThenShouldReceiveAnError(t *testing.T) {
	task := &entity.Task{
		TaskId:      "id",
		Name:        "",
		Description: "desc",
		Config:      map[string]interface{}{},
		Schedule:    entity.Schedule{Expression: "0 * * * *"},
		Status:      entity.TaskStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := task.IsValid()
	assert.EqualError(t, err, "invalid task ID")
}

func TestGivenAnEmptyDescription_WhenCreateANewTask_ThenShouldReceiveAnError(t *testing.T) {
	task := &entity.Task{
		TaskId:      "id",
		Name:        "name",
		Description: "",
		Config:      map[string]interface{}{},
		Schedule:    entity.Schedule{Expression: "0 * * * *"},
		Status:      entity.TaskStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := task.IsValid()
	assert.EqualError(t, err, "invalid task description")
}

func TestGivenNilConfig_WhenCreateANewTask_ThenShouldReceiveAnError(t *testing.T) {
	task := &entity.Task{
		TaskId:      "id",
		Name:        "name",
		Description: "desc",
		Config:      nil,
		Schedule:    entity.Schedule{Expression: "0 * * * *"},
		Status:      entity.TaskStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := task.IsValid()
	assert.EqualError(t, err, "invalid task config")
}

func TestGivenInvalidStatus_WhenCreateANewTask_ThenShouldReceiveAnError(t *testing.T) {
	task := &entity.Task{
		TaskId:      "id",
		Name:        "name",
		Description: "desc",
		Config:      map[string]interface{}{},
		Schedule:    entity.Schedule{Expression: "0 * * * *"},
		Status:      "unknown",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := task.IsValid()
	assert.EqualError(t, err, "invalid task status")
}

func TestGivenInvalidCron_WhenCreateANewTask_ThenShouldReceiveAnError(t *testing.T) {
	task := &entity.Task{
		TaskId:      "id",
		Name:        "name",
		Description: "desc",
		Config:      map[string]interface{}{},
		Schedule:    entity.Schedule{Expression: "invalid cron"},
		Status:      entity.TaskStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := task.IsValid()
	assert.Error(t, err)
}

func TestGivenValidParams_WhenCallNewTask_ThenShouldCreateTaskWithAllParams(t *testing.T) {
	task, err := callNewTask("task-1", "Test Task", "A test task", map[string]interface{}{"foo": "bar"}, "0 * * * *")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, entity.TaskID("task-1"), task.TaskId)
	assert.Equal(t, "Test Task", task.Name)
	assert.Equal(t, "A test task", task.Description)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, task.Config)
	assert.Equal(t, "0 * * * *", task.Schedule.Expression)
	assert.Equal(t, entity.TaskStatusActive, task.Status)
}

func TestTaskExecutionFields(t *testing.T) {
	now := time.Now()
	completed := now.Add(1 * time.Minute)
	result := map[string]interface{}{"output": "ok"}
	exec := entity.TaskExecution{
		TaskExecutionId: "exec-1",
		TaskID:          "task-1",
		WorkerID:        "worker-1",
		StartedAt:       now,
		CompletedAt:     &completed,
		Success:         true,
		Error:           "",
		Result:          result,
		RetryCount:      1,
		Duration:        completed.Sub(now),
	}
	assert.Equal(t, "exec-1", exec.TaskExecutionId)
	assert.Equal(t, "task-1", string(exec.TaskID))
	assert.Equal(t, "worker-1", exec.WorkerID)
	assert.Equal(t, now, exec.StartedAt)
	assert.Equal(t, &completed, exec.CompletedAt)
	assert.True(t, exec.Success)
	assert.Equal(t, "", exec.Error)
	assert.Equal(t, result, exec.Result)
	assert.Equal(t, 1, exec.RetryCount)
	assert.Equal(t, 1*time.Minute, exec.Duration)
}

// Helpers to call unexported logic from entity

func callNewTask(id entity.TaskID, name, desc string, config map[string]interface{}, expr string) (*entity.Task, error) {
	task := &entity.Task{
		TaskId:      id,
		Name:        name,
		Description: desc,
		Config:      config,
		Schedule: entity.Schedule{
			Expression: expr,
			NextRun:    nil,
			LastRun:    nil,
		},
		Status:    entity.TaskStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := task.IsValid()
	if err != nil {
		return nil, err
	}
	return task, nil
}
