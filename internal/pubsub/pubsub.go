package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	newChannel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	switch queueType {
	case SimpleQueueDurable:
		newQueue, err := newChannel.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
		err = newChannel.QueueBind(queueName, key, exchange, false, nil)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
		return newChannel, newQueue, nil
	case SimpleQueueTransient:
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
