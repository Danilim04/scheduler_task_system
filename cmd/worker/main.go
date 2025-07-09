package main

import (
	"context"
	"log"
	"scheduler_task_system/internal/app"
	messagingengine "scheduler_task_system/internal/infra/menssaging_enginer"
)

func main() {
	ctx := context.Background()

	channel, conn, err := messagingengine.OpenChannel()
	if err != nil {
		panic(err)
	}

	engine, err := messagingengine.NewMessagingEngine(channel)
	defer messagingengine.Disconnect(engine, conn)

	if err != nil {
		log.Fatalf("engine: %v", err)
	}

	reg := app.Default

	worker := app.NewWorker(engine.Delivery, reg)

	log.Println("Worker started")
	worker.Start(ctx)
}
