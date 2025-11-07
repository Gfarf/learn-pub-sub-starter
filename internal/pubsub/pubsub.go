package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	durable SimpleQueueType = iota
	transient
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: buf})
	return err
}

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	newChannel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	switch queueType {
	case durable:
		newQueue, err := newChannel.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
		err = newChannel.QueueBind(queueName, key, exchange, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
		return newChannel, newQueue, nil
	case transient:
		newQueue, err := newChannel.QueueDeclare(queueName, false, true, true, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
		err = newChannel.QueueBind(queueName, key, exchange, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
		return newChannel, newQueue, nil
	default:
		return nil, amqp.Queue{}, fmt.Errorf("queue type not identified")
	}
}
