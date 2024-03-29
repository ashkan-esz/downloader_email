package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ICreateQueue is the interface for creating
type ICreateQueue interface {
	CreateQueue(config ConfigQueue) (queue amqp.Queue, err error)
}

// IQueueBinder is the interface for binding and unbinding queues
type IQueueBinder interface {
	BindQueueExchange(config ConfigBindQueue) (err error)
	UnbindQueueExchange(config ConfigBindQueue) (err error)
}

const (
	EmailQueue      = "email"
	EmailBindingKey = "email"
)

func (r *rabbit) createQueuesAndBind() {
	emailConfig := ConfigQueue{
		Name:       EmailQueue,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	_, err := r.CreateQueue(emailConfig)
	if err != nil {
		log.Printf("error creating queue %s: %s\n", EmailQueue, err)
	}

	emailConfigBindConfig := ConfigBindQueue{
		QueueName:  EmailQueue,
		Exchange:   EmailExchange,
		RoutingKey: EmailBindingKey,
		NoWait:     false,
	}
	err = r.BindQueueExchange(emailConfigBindConfig)
	if err != nil {
		log.Printf("error binding queue %s: %s\n", EmailQueue, err)
	}
}

// CreateQueue creates a queue
func (r *rabbit) CreateQueue(config ConfigQueue) (queue amqp.Queue, err error) {
	if r.chConsumer == nil {
		err = amqp.ErrClosed
		return
	}
	queue, err = r.chConsumer.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)
	return
}

// BindQueueExchange binds a queue to an exchange
func (r *rabbit) BindQueueExchange(config ConfigBindQueue) (err error) {
	if r.chConsumer == nil {
		err = amqp.ErrClosed
		return
	}
	err = r.chConsumer.QueueBind(
		config.QueueName,
		config.RoutingKey,
		config.Exchange,
		config.NoWait,
		config.Args,
	)
	return
}

// UnbindQueueExchange unbinds a queue from an exchange
func (r *rabbit) UnbindQueueExchange(config ConfigBindQueue) (err error) {
	if r.chConsumer == nil {
		err = amqp.ErrClosed
		return
	}
	err = r.chConsumer.QueueUnbind(
		config.QueueName,
		config.RoutingKey,
		config.Exchange,
		config.Args,
	)
	return
}
