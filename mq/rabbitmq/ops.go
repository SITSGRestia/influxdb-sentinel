package rabbitmq

import (
	"github.com/streadway/amqp"
	"fmt"
	"errors"
)

type Ops struct {
	url        string
	kind       string
	vHost      string
	exchange   string
	routingKey string
}

func (ops *Ops) Send(msg string) error {
	conn, err := amqp.Dial(ops.url)
	if err != nil {
		return errors.New(fmt.Sprintf("RabbitMQ连接失败.%+v", err))
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return errors.New(fmt.Sprintf("RabbitMQ管道打开失败.%+v", err))
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(
		ops.exchange, // name
		ops.kind,     // type
		false,        // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return errors.New(fmt.Sprintf("RabbitMQ的exchange定义失败.%+v", err))
	}
	err = ch.Publish(
		ops.exchange,   // exchange
		ops.routingKey, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Expiration:  "60000",
			Body:        []byte(msg),
		})
	if err != nil {
		return errors.New(fmt.Sprintf("RabbitMQ发送消息失败.%+v", err))
	}
	return nil
}
