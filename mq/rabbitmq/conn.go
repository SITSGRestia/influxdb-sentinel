package rabbitmq

import (
	"fmt"
)

type RabbitMQ struct {
	Host       string `json:"host" description:"消息队列地址"`
	Port       int    `json:"port" description:"消息队列端口"`
	Username   string `json:"username" description:"消息队列访问用户名"`
	Password   string `json:"password" description:"消息队列访问密码"`
	Kind       string `json:"kind" description:"消息队列访问密码"`
	VHost      string `json:"vHost" description:"消息队列VHost名"`
	Exchange   string `json:"exchange" description:"消息队列Exchange名"`
	RoutingKey string `json:"routingKey" description:"消息队列名"`
}

func New(option RabbitMQ) *Ops {
	return &Ops{url: fmt.Sprintf("amqp://%s:%s@%s:%d/%s", option.Username, option.Password, option.Host, option.Port, option.VHost), kind: option.Kind, vHost: option.VHost, exchange: option.Exchange, routingKey: option.RoutingKey}
}
