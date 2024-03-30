package rabbitmq

import (
	"downloader_email/pkg"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ICreateExchange interface {
	CreateExchange(config ConfigExchange) (err error)
}

const (
	EmailExchange     = "EmailExchange"
	EmailExchangeType = "direct"
)

func (r *rabbit) createExchanges() {
	emailConfig := ConfigExchange{
		Name:       EmailExchange,
		Type:       EmailExchangeType,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
	err := r.CreateExchange(emailConfig)
	if err != nil {
		message := fmt.Sprintf("error creating exchange %v: %s", EmailExchange, err)
		pkg.SaveError(message, err)
	}
}

// CreateExchange creates an exchange
func (r *rabbit) CreateExchange(config ConfigExchange) (err error) {
	if r.chConsumer == nil {
		return amqp.ErrClosed
	}
	err = r.chConsumer.ExchangeDeclare(
		config.Name,
		config.Type,
		config.Durable,
		config.AutoDelete,
		config.Internal,
		config.NoWait,
		config.Args,
	)
	return
}
