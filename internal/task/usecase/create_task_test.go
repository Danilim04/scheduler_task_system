package usecase_test

import (
	"os"
	"scheduler_task_system/internal/task/infra/codetemplate"
	"scheduler_task_system/internal/task/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	repository, err := codetemplate.NewTaskTemplateRepository(os.Getenv("GO_ROOTPATH"))
	assert.Nil(t, repository)
	assert.Error(t, err)
	uc := usecase.NewCreateTaskUseCase(repository)

}
