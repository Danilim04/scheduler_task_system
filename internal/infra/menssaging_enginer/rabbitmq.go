package messagingengine

import (
	"context"
	"scheduler_task_system/internal/core/entity"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessagingEngine struct {
	Channel  *amqp.Channel
	Delivery <-chan amqp.Delivery
}

func NewMessagingEngine(ch *amqp.Channel) (*MessagingEngine, error) {
	deliveries, err := ch.Consume(
		"tasks",
		"go-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &MessagingEngine{
		Channel:  ch,
		Delivery: deliveries,
	}, nil
}

func (m *MessagingEngine) Publish(ctx context.Context, task entity.TaskExecution, payload []byte) error {
	return m.Channel.PublishWithContext(
		ctx,
		"amq.direct",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			Headers: amqp.Table{
				"task_id": task.TaskID,
			},
		},
	)
}

func OpenChannel() (*amqp.Channel, *amqp.Connection, error) {
	// 1) Conexão
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, nil, err
	}

	// 2) Canal
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close() // cleanup antes do return
		return nil, nil, err
	}

	// 3) QoS
	if err := ch.Qos(100, 0, false); err != nil {
		_ = ch.Close() // cleanup
		_ = conn.Close()
		return nil, nil, err
	}

	// 4) Sucesso
	return ch, conn, nil
}

func Disconnect(e *MessagingEngine, conn *amqp.Connection) error {
	_ = e.Channel.Cancel("go-consumer", false)
	_ = e.Channel.Close()
	return conn.Close()
}
