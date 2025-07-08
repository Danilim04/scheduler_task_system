package app

import (
	"errors"
	"scheduler_task_system/internal/core/port"
)

type Factory func() port.TaskExecuter

type Registry struct {
	factories map[string]Factory
}

func NewRegistry() *Registry {
	return &Registry{factories: make(map[string]Factory)}
}

func (r *Registry) Register(name string, taskExecuter Factory) {
	r.factories[name] = taskExecuter
}

func (r *Registry) Resolve(name string) (port.TaskExecuter, error) {
	f, ok := r.factories[name]
	if !ok {
		return nil, errors.New("task não existe ou não foi registrada")
	}
	return f(), nil
}

var Default = NewRegistry()
