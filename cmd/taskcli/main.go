package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"scheduler_task_system/internal/core/usecase"
	"scheduler_task_system/internal/infra/mongodb"
	"scheduler_task_system/internal/infra/template"
)

func main() {
	// ----------- Flags -----------
	taskID := flag.String("id", "", "ID único da task (obrigatório)")
	name := flag.String("name", "", "Nome da task (obrigatório)")
	desc := flag.String("desc", "", "Descrição da task")
	cronExpr := flag.String("cron", "0 * * * *", "Expressão CRON")
	payload := flag.String("payload", "{}", "JSON bruto do payload")
	rootPath := flag.String("root", os.Getenv("GO_ROOTPATH"), "Path dos templates")

	flag.Parse()

	if *taskID == "" || *name == "" {
		flag.Usage()
		log.Fatal("as flags -id e -name são obrigatórias")
	}

	payloadBytes := []byte(*payload)
	if !json.Valid(payloadBytes) {
		log.Fatalf("payload não é JSON válido: %s", *payload)
	}

	// ----------- Boot de infra -----------
	ctx := context.Background()

	repoTpl, err := template.NewTaskTemplateRepository(*rootPath)
	if err != nil {
		log.Fatal(err)
	}

	client, err := mongodb.ConnectMongodb()
	if err != nil {
		log.Fatal(err)
	}
	defer mongodb.DisconnectMongodb(ctx, client)

	repoMongo := mongodb.NewTaskRepositoryMongo(client)

	// ----------- Caso de uso -----------
	uc := usecase.NewCreateTaskUseCase(repoMongo, repoTpl)

	input := usecase.CreateTaskInputDto{
		TaskId:      *taskID,
		Name:        *name,
		Description: *desc,
		Payload:     payloadBytes,
		Expression:  *cronExpr,
	}

	exec, err := uc.Execute(ctx, input)
	if err != nil {
		log.Fatalf("falha ao criar task: %v", err)
	}
	fmt.Println("Task criada com sucesso:", exec)
}
