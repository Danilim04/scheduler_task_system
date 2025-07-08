package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"scheduler_task_system/internal/core/usecase"
	"scheduler_task_system/internal/infra/mongodb"
	"scheduler_task_system/internal/infra/template"
)

var ctx context.Context

var rootpath string = os.Getenv("GO_ROOTPATH")

func main() {
	var err error
	repositoryTemplate, err := template.NewTaskTemplateRepository(rootpath)
	if err != nil {
		panic(err)
	}
	client, err := mongodb.ConnectMongodb()
	if err != nil {
		panic(err)
	}
	defer mongodb.DisconnectMongodb(ctx, client)
	repositoryMongo, err := mongodb.NewTaskRepositoryMongo(client)
	if err != nil {
		log.Fatal(err)
	}

	uc := usecase.NewCreateTaskUseCase(repositoryMongo, repositoryTemplate)

	payloadBytes, err := json.Marshal(map[string]interface{}{
		"key": "value",
	})
	if err != nil {
		log.Fatal(err)
	}

	input := usecase.CreateTaskInputDto{
		TaskId:      "task_test",
		Name:        "task_test",
		Description: "Test description",
		Payload:     payloadBytes,
		Expression:  "0 * * * *",
	}

	exec, err := uc.Execute(ctx, input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(exec)

}
