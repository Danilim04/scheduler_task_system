package app

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const headerKey = "task" // mesmo nome usado no Publish

type Worker struct {
	deliveries <-chan amqp.Delivery
	registry   *Registry
}

func NewWorker(deliveries <-chan amqp.Delivery, r *Registry) *Worker {
	return &Worker{deliveries: deliveries, registry: r}
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-w.deliveries:
			if !ok {
				return
			}

			raw, ok := msg.Headers[headerKey]
			if !ok {
				log.Printf("msg sem header %q -> reject", headerKey)
				_ = msg.Reject(false)
				continue
			}
			taskName, ok := raw.(string)
			if !ok {
				log.Printf("header %q não é string -> reject", headerKey)
				_ = msg.Reject(false)
				continue
			}

			executer, err := w.registry.Resolve(taskName)
			if err != nil {
				log.Printf("task %q não registrada: %v -> reject", taskName, err)
				_ = msg.Reject(false)
				continue
			}

			if err := executer.Execute(ctx, msg.Body); err != nil {
				log.Printf("[%s] erro: %v -> nack requeue", taskName, err)
				_ = msg.Nack(false, true) // requeue
				continue
			}

			_ = msg.Ack(false)
		}
	}
}
