package usecase_test

import (
	"context"
	"errors"
	"testing"

	"scheduler_task_system/internal/task/entity"
	"scheduler_task_system/internal/task/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do TaskRepositoryInterface para testes
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) ExistsByID(ctx context.Context, id entity.TaskID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockTaskRepository) Save(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id entity.TaskID) (*entity.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAll(ctx context.Context) ([]*entity.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) DeleteByID(ctx context.Context, id entity.TaskID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

type MockTaskRepositoryTemplate struct {
	mock.Mock
}

func (m *MockTaskRepositoryTemplate) CreateTemplate(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func TestGivenEmptyName_WhenExecuteCreateTask_ThenShouldReceiveError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskRepository)
	mockTemplate := new(MockTaskRepositoryTemplate)
	uc := usecase.NewCreateTaskUseCase(mockRepo, mockTemplate)

	input := usecase.CreateTaskInputDto{
		TaskId:      "task_test",
		Name:        "", // nome vazio
		Description: "Test description",
		Config: map[string]interface{}{
			"key": "value",
		},
		Expression: "0 * * * *",
	}

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid task name")
	mockTemplate.AssertNotCalled(t, "CreateTemplate")
}

func TestGivenEmptyDescription_WhenExecuteCreateTask_ThenShouldReceiveError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskRepository)
	mockTemplate := new(MockTaskRepositoryTemplate)
	uc := usecase.NewCreateTaskUseCase(mockRepo, mockTemplate)

	input := usecase.CreateTaskInputDto{
		TaskId:      "task_test",
		Name:        "Test Task",
		Description: "", // descrição vazia
		Config: map[string]interface{}{
			"key": "value",
		},
		Expression: "0 * * * *",
	}

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid task description")
	mockTemplate.AssertNotCalled(t, "CreateTemplate")
}

func TestGivenNilConfig_WhenExecuteCreateTask_ThenShouldReceiveError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskRepository)
	mockTemplate := new(MockTaskRepositoryTemplate)
	uc := usecase.NewCreateTaskUseCase(mockRepo, mockTemplate)

	input := usecase.CreateTaskInputDto{
		TaskId:      "task_test",
		Name:        "Test Task",
		Description: "Test description",
		Config:      nil, // config nil
		Expression:  "0 * * * *",
	}

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid task config")
	mockTemplate.AssertNotCalled(t, "CreateTemplate")
}

func TestGivenEmptyConfig_WhenExecuteCreateTask_ThenShouldReceiveError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskRepository)
	mockTemplate := new(MockTaskRepositoryTemplate)
	uc := usecase.NewCreateTaskUseCase(mockRepo, mockTemplate)

	input := usecase.CreateTaskInputDto{
		TaskId:      "task_test",
		Name:        "Test Task",
		Description: "Test description",
		Config:      map[string]interface{}{}, // config vazio
		Expression:  "0 * * * *",
	}

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid task config")
	mockTemplate.AssertNotCalled(t, "CreateTemplate")
}

func TestGivenInvalidCron_WhenExecuteCreateTask_ThenShouldReceiveError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskRepository)
	mockTemplate := new(MockTaskRepositoryTemplate)
	uc := usecase.NewCreateTaskUseCase(mockRepo, mockTemplate)

	input := usecase.CreateTaskInputDto{
		TaskId:      "task_test",
		Name:        "Test Task",
		Description: "Test description",
		Config: map[string]interface{}{
			"key": "value",
		},
		Expression: "invalid cron", // cron inválido
	}

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Nil(t, result)
	assert.Error(t, err)
	mockTemplate.AssertNotCalled(t, "CreateTemplate")
}

func TestGivenTemplateRepositoryError_WhenExecuteCreateTask_ThenShouldReceiveError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskRepository)
	mockTemplate := new(MockTaskRepositoryTemplate)
	mockTemplate.On("CreateTemplate", mock.Anything, mock.Anything).Return(errors.New("template creation error"))

	uc := usecase.NewCreateTaskUseCase(mockRepo, mockTemplate)

	input := usecase.CreateTaskInputDto{
		TaskId:      "task_test",
		Name:        "Test Task",
		Description: "Test description",
		Config: map[string]interface{}{
			"key": "value",
		},
		Expression: "0 * * * *",
	}

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Nil(t, result)
	assert.EqualError(t, err, "template creation error")
	mockTemplate.AssertExpectations(t)
}

func TestGivenValidInput_WhenExecuteCreateTask_ThenShouldCreateTaskSuccessfully(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskRepository)
	mockTemplate := new(MockTaskRepositoryTemplate)
	mockTemplate.On("CreateTemplate", mock.Anything, mock.Anything).Return(nil)

	uc := usecase.NewCreateTaskUseCase(mockRepo, mockTemplate)

	input := usecase.CreateTaskInputDto{
		TaskId:      "task-1",
		Name:        "Test Task",
		Description: "A test task",
		Config: map[string]interface{}{
			"foo": "bar",
		},
		Expression: "0 * * * *",
	}

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.TaskID("task-1"), result.TaskId)
	assert.Equal(t, "Test Task", result.Name)
	assert.Equal(t, "A test task", result.Description)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, result.Config)
	assert.Equal(t, "0 * * * *", result.Schedule.Expression)
	assert.Equal(t, entity.TaskStatusActive, result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
	assert.Nil(t, result.Schedule.NextRun)
	assert.Nil(t, result.Schedule.LastRun)

	mockTemplate.AssertExpectations(t)
}

func TestGivenValidInputWithDifferentCronExpressions_WhenExecuteCreateTask_ThenShouldCreateTaskSuccessfully(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskRepository)
	mockTemplate := new(MockTaskRepositoryTemplate)

	uc := usecase.NewCreateTaskUseCase(mockRepo, mockTemplate)

	testCases := []struct {
		name       string
		expression string
	}{
		{"hourly", "0 * * * *"},
		{"daily", "0 0 * * *"},
		{"weekly", "0 0 * * 0"},
		{"monthly", "0 0 1 * *"},
		{"every_15_minutes", "*/15 * * * *"},
		{"every_30_minutes", "*/30 * * * *"},
		{"every_hour_at_30", "30 * * * *"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configurar mock para cada test case
			mockTemplate.On("CreateTemplate", mock.Anything, mock.Anything).Return(nil).Once()

			input := usecase.CreateTaskInputDto{
				TaskId:      "task-" + tc.name,
				Name:        "Test Task " + tc.name,
				Description: "A test task for " + tc.name,
				Config: map[string]interface{}{
					"type": tc.name,
				},
				Expression: tc.expression,
			}

			// Act
			result, err := uc.Execute(context.Background(), input)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.expression, result.Schedule.Expression)
			assert.Equal(t, entity.TaskID("task-"+tc.name), result.TaskId)
			assert.Equal(t, "Test Task "+tc.name, result.Name)
		})
	}

	mockTemplate.AssertExpectations(t)
}

func TestGivenDifferentConfigTypes_WhenExecuteCreateTask_ThenShouldCreateTaskSuccessfully(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskRepository)
	mockTemplate := new(MockTaskRepositoryTemplate)

	uc := usecase.NewCreateTaskUseCase(mockRepo, mockTemplate)

	testCases := []struct {
		name   string
		config map[string]interface{}
	}{
		{
			"string_config",
			map[string]interface{}{
				"url":    "https://example.com",
				"method": "GET",
			},
		},
		{
			"number_config",
			map[string]interface{}{
				"timeout": 30,
				"retries": 3,
			},
		},
		{
			"boolean_config",
			map[string]interface{}{
				"enabled": true,
				"async":   false,
			},
		},
		{
			"mixed_config",
			map[string]interface{}{
				"name":   "test",
				"count":  10,
				"active": true,
				"tags":   []string{"tag1", "tag2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configurar mock para cada test case
			mockTemplate.On("CreateTemplate", mock.Anything, mock.Anything).Return(nil).Once()

			input := usecase.CreateTaskInputDto{
				TaskId:      "task-" + tc.name,
				Name:        "Test Task " + tc.name,
				Description: "A test task for " + tc.name,
				Config:      tc.config,
				Expression:  "0 * * * *",
			}

			// Act
			result, err := uc.Execute(context.Background(), input)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.config, result.Config)
			assert.Equal(t, entity.TaskID("task-"+tc.name), result.TaskId)
		})
	}

	mockTemplate.AssertExpectations(t)
}
