package template

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"scheduler_task_system/internal/core/entity"
	"strings"
	"text/template"
)

type TaskRepositoryTemplate struct {
	rootPath          string
	TaskCodeTemplates *TaskCodeTemplates
}
type TaskCodeTemplates struct {
	TaskCodeTemplateExecuter *template.Template
}

func NewTaskTemplateRepository(rootPath string) (*TaskRepositoryTemplate, error) {

	TaskCodeTemplatesDir := map[string]string{
		"executer": rootPath + "/internal/infra/template/executer_template.tmpl",
	}
	tt, err := LoadTemplates(TaskCodeTemplatesDir)
	if err != nil {
		return nil, errors.New("erro ao carregar os templates de task" + err.Error())
	}

	tr := &TaskRepositoryTemplate{
		rootPath:          rootPath,
		TaskCodeTemplates: tt,
	}

	return tr, nil
}

func LoadTemplates(tr map[string]string) (*TaskCodeTemplates, error) {
	var err error

	te, err := template.ParseFiles(tr["executer"])
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar template entity: %v", err)
	}
	templates := &TaskCodeTemplates{
		TaskCodeTemplateExecuter: te,
	}

	return templates, nil
}

func (tr *TaskRepositoryTemplate) Generate(ctx context.Context, task *entity.Task) error {

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

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return errors.New("erro ao criar diretório" + dir + ": " + err.Error())
		}
		if err := os.Chown(dir, 1000, 1000); err != nil {
			return fmt.Errorf("erro ao definir propriedade do diretório %s: %v", dir, err)
		}
	}

	file, err := os.Create(taskDir + "/executer.go")
	if err != nil {
		return err
	}
	if err := tr.TaskCodeTemplates.TaskCodeTemplateExecuter.Execute(file, map[string]interface{}{"TaskName": taskName}); err != nil {
		return fmt.Errorf("erro ao executar template: %v", err)
	}
	if err := os.Chown(taskDir+"/executer.go", 1000, 1000); err != nil {
		return fmt.Errorf("erro ao definir propriedade do diretório %s: %v", taskDir, err)
	}
	return nil
}
