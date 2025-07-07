package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"scheduler_task_system/internal/core/usecase"
	"scheduler_task_system/internal/infra/mongodb"
	"scheduler_task_system/internal/infra/template"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ctx context.Context

var rootpath string = os.Getenv("GO_ROOTPATH")

func main() {
	var err error
	repositoryTemplate, err := template.NewTaskTemplateRepository(rootpath)
	if err != nil {
		panic(err)
	}
	client, err := connectMongodb()
	if err != nil {
		panic(err)
	}
	defer disconnectMongodb(client)
	repositoryMongo := mongodb.NewTaskRepositoryMongo(client)
	uc := usecase.NewCreateTaskUseCase(repositoryMongo, repositoryTemplate)

	input := usecase.CreateTaskInputDto{
		TaskId:      "task_test",
		Name:        "task_test",
		Description: "Test description",
		Config: map[string]interface{}{
			"key": "value",
		},
		Expression: "0 * * * *",
	}

	exec, err := uc.Execute(ctx, input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(exec)

}

func connectMongodb() (*mongo.Client, error) {
	uri := "mongodb://dev:dev123@mongodb:27017/"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, errors.New("falha ao conectar no banco de dados: " + err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.New("falha ao verificar conexão com o banco: " + err.Error())
	}

	return client, nil
}

func disconnectMongodb(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Disconnect(ctx)
}
