package entity_test

import (
	"fmt"
	"testing"

	"scheduler_task_system/internal/task/entity"

	"github.com/stretchr/testify/assert"
)

func TestGivenAnEmptyName_WhenCreateANewTask_ThenShouldReceiveAnError(t *testing.T) {
	task, err := entity.NewCreateTask(
		"task_test",
		"",
		"desc",
		map[string]interface{}{
			"key": "value",
		},
		"0 * * * *",
	)
	assert.Nil(t, task)
	assert.EqualError(t, err, "erro ao carregar os templates de task")
	fmt.Println(task, err)
}

func TestGivenAnEmptyDescription_WhenCreateANewTask_ThenShouldReceiveAnError(t *testing.T) {
	task, err := entity.NewCreateTask(
		"task_test",
		"Test Task",
		"", // descrição vazia
		map[string]interface{}{
			"key": "value",
		},
		"0 * * * *",
	)
	assert.Nil(t, task)
	assert.EqualError(t, err, "invalid task description")
}

func TestGivenNilConfig_WhenCreateANewTask_ThenShouldReceiveAnError(t *testing.T) {
	task, err := entity.NewCreateTask(
		"task_test",
		"Test Task",
		"Test description",
		nil, // config nil
		"0 * * * *",
	)
	assert.Nil(t, task)
	assert.EqualError(t, err, "invalid task config")
}

func TestGivenInvalidCron_WhenCreateANewTask_ThenShouldReceiveAnError(t *testing.T) {
	task, err := entity.NewCreateTask(
		"task_test",
		"Test Task",
		"Test description",
		map[string]interface{}{
			"key": "value",
		},
		"invalid cron", // cron inválido
	)
	assert.Nil(t, task)
	assert.Error(t, err)
}

func TestGivenValidParams_WhenCallNewTask_ThenShouldCreateTaskWithAllParams(t *testing.T) {
	task, err := entity.NewCreateTask(
		"task-1",
		"Test Task",
		"A test task",
		map[string]interface{}{"foo": "bar"},
		"0 * * * *",
	)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, entity.TaskID("task-1"), task.TaskId)
	assert.Equal(t, "Test Task", task.Name)
	assert.Equal(t, "A test task", task.Description)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, task.Config)
	assert.Equal(t, "0 * * * *", task.Schedule.Expression)
	assert.Equal(t, entity.TaskStatusActive, task.Status)
}
