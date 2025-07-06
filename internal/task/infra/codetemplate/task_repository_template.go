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
	rootPath          string
	TaskCodeTemplates *TaskCodeTemplates
}
type TaskCodeTemplates struct {
	TaskCodeTemplatesEntity  *template.Template
	TaskCodeTemplatesUseCase *template.Template
	TaskCodeTemplatesInfra   *template.Template
}

func NewTaskTemplateRepository(rootPath string) (*TaskRepositoryTemplate, error) {

	TaskCodeTemplatesDir := map[string]string{
		"entity":  rootPath + "/internal/task/taskCodeTemplate/entity/taskEntityTemplate.tmpl",
		"useCase": rootPath + "/internal/task/taskCodeTemplate/entity/taskEntityTemplate.tmpl",
		"infra":   rootPath + "/internal/task/taskCodeTemplate/entity/taskEntityTemplate.tmpl",
	}
	tt, err := LoadTemplates(TaskCodeTemplatesDir)
	if err != nil {
		return nil, errors.New("erro ao carregar os templates de task")
	}

	tr := &TaskRepositoryTemplate{
		rootPath:          rootPath,
		TaskCodeTemplates: tt,
	}

	return tr, nil
}

func LoadTemplates(tr map[string]string) (*TaskCodeTemplates, error) {
	var err error

	te, err := template.ParseFiles(tr["entity"])
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar template entity: %v", err)
	}

	tu, err := template.ParseFiles(tr["useCase"])
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar template useCase: %v", err)
	}

	ti, err := template.ParseFiles(tr["infra"])
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar template infra: %v", err)
	}

	templates := &TaskCodeTemplates{
		TaskCodeTemplatesEntity:  te,
		TaskCodeTemplatesUseCase: tu,
		TaskCodeTemplatesInfra:   ti,
	}

	return templates, nil
}

func (tr *TaskRepositoryTemplate) CreateTemplate(ctx context.Context, task *entity.Task) error {

	taskName := task.Name

	taskDir := filepath.Join(tr.rootPath, "internal", strings.ToLower(taskName))

	dirs := []string{
		filepath.Join(taskDir, "entity"),
		filepath.Join(taskDir, "useCase"),
		filepath.Join(taskDir, "infra"),
	}

	if err := os.MkdirAll(taskDir, 0777); err != nil {
		return errors.New("erro ao criar diretório" + taskDir + ": " + err.Error())
	}
	if err := os.Chown(taskDir, 1000, 1000); err != nil {
		return fmt.Errorf("erro ao definir propriedade do diretório %s: %v", taskDir, err)
	}

	var dirsName []string

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return errors.New("erro ao criar diretório" + dir + ": " + err.Error())
		}
		if err := os.Chown(dir, 1000, 1000); err != nil {
			return fmt.Errorf("erro ao definir propriedade do diretório %s: %v", dir, err)
		}
		dirsName = append(dirsName, dir+"/"+taskName+".go")
	}

	for _, dirName := range dirsName {
		file, err := os.Create(dirName)
		if err := os.Chown(dirName, 1000, 1000); err != nil {
			return fmt.Errorf("erro ao definir propriedade do diretório %s: %v", dirName, err)
		}
		if err != nil {
			return errors.New("erro ao criar arquivo")
		}
		defer file.Close()

		if err := tr.TaskCodeTemplates.TaskCodeTemplatesEntity.Execute(file, map[string]interface{}{"EntityName": taskName}); err != nil {
			return fmt.Errorf("erro ao executar template: %v", err)
		}
		if err := tr.TaskCodeTemplates.TaskCodeTemplatesUseCase.Execute(file, map[string]interface{}{"EntityName": taskName}); err != nil {
			return fmt.Errorf("erro ao executar template: %v", err)
		}
		if err := tr.TaskCodeTemplates.TaskCodeTemplatesInfra.Execute(file, map[string]interface{}{"EntityName": taskName}); err != nil {
			return fmt.Errorf("erro ao executar template: %v", err)
		}
	}

	return nil
}
