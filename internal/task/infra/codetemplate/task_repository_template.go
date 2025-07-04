package codetemplate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"scheduler_task_system/internal/task/entity"
	"strings"
	"text/template"
)

type TaskRepositoryTemplate struct {
	rootPath             string
	TaskCodeTemplatesDir map[string]string
}

func NewTaskTemplateRepository(rootPath string) (*TaskRepositoryTemplate, error) {
	tr := &TaskRepositoryTemplate{
		rootPath: rootPath,
		TaskCodeTemplatesDir: map[string]string{
			"entity":  rootPath + "/internal/task/taskCodeTemplate/taskEntityTemplate",
			"useCase": rootPath + "/internal/task/taskCodeTemplate/taskEntityTemplate",
			"infra":   rootPath + "/internal/task/taskCodeTemplate/taskEntityTemplate",
		},
	}
	if err := LoadTemplates(tr); err != nil {
		return nil, errors.New("erro ao carregar os templates de task")
	}
	return tr, nil
}

var TaskCodeTemplates = map[string]*template.Template{}

func LoadTemplates(tr *TaskRepositoryTemplate) error {
	var err error
	TaskCodeTemplates["entity"], err = template.ParseFiles(tr.TaskCodeTemplatesDir["entity"])
	if err != nil {
		return err
	}
	TaskCodeTemplates["useCase"], err = template.ParseFiles(tr.TaskCodeTemplatesDir["useCase"])
	if err != nil {
		return err
	}
	TaskCodeTemplates["infra"], err = template.ParseFiles(tr.TaskCodeTemplatesDir["useCase"])
	if err != nil {
		return err
	}
	return nil
}

func (tr *TaskRepositoryTemplate) CreateTemplate(ctx context.Context, task *entity.Task) error {

	taskName := task.Name

	taskDir := filepath.Join(tr.rootPath, "internal", strings.ToLower(taskName))

	dirs := []string{
		filepath.Join(taskDir, "entity"),
		filepath.Join(taskDir, "useCase"),
		filepath.Join(taskDir, "infra"),
	}

	var dirsName []string

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.New("erro ao criar diretório" + dir + ": " + err.Error())
		}
		dirsName = append(dirsName, dir+taskName)
	}

	for _, dirName := range dirsName {
		file, err := os.Create(dirName)
		if err != nil {
			return errors.New("erro ao criar arquivo")
		}
		defer file.Close()
		if err := TaskCodeTemplates["entity"].Execute(file, map[string]interface{}{"EntityName": taskName}); err != nil {
			return fmt.Errorf("erro ao executar template: %v", err)
		}
		if err := TaskCodeTemplates["useCase"].Execute(file, map[string]interface{}{"EntityName": taskName}); err != nil {
			return fmt.Errorf("erro ao executar template: %v", err)
		}
		if err := TaskCodeTemplates["infra"].Execute(file, map[string]interface{}{"EntityName": taskName}); err != nil {
			return fmt.Errorf("erro ao executar template: %v", err)
		}
	}

	return nil
}
