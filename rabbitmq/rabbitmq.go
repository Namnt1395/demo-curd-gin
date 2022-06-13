package rabbitmq

import (
	"demo-curd/config"
	"demo-curd/util"
	"demo-curd/util/constant"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection      *amqp.Connection
	Channel         *amqp.Channel
	ProdExchange    *string
	ProdRoutingKeys map[string]string
	Queue           *amqp.Queue
}

func NewRabbitMQ(c config.Config) (*RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v", c.RabbitMQ.Username, c.RabbitMQ.Password, c.RabbitMQ.Host, c.RabbitMQ.Port))
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// config producer
	prodRks := make(map[string]string, 0)
	var exchange *string
	if c.RabbitMQ.Producer != nil {
		exchange = &c.RabbitMQ.Producer.Exchange
		err = channel.ExchangeDeclare(
			c.RabbitMQ.Producer.Exchange, // name
			"topic",                      // type
			true,                         // durable
			false,                        // auto-deleted
			false,                        // internal
			false,                        // no-wait
			nil,                          // arguments
		)
		if err != nil {
			return nil, err
		}

		if c.RabbitMQ.Producer.Bindings != nil {
			for key, binding := range c.RabbitMQ.Producer.Bindings {
				queue, err := channel.QueueDeclare(
					binding.Queue, // name
					true,          // durable
					false,         // delete when unused
					false,         // exclusive
					false,         // no-wait
					nil,           // arguments
				)
				if err != nil {
					return nil, err
				}

				routingKey := ""
				if binding.RoutingKey != nil {
					routingKey = fmt.Sprintf("%v-%v", c.RabbitMQ.Producer.Exchange, *binding.RoutingKey)
					prodRks[key] = routingKey
				}
				err = channel.QueueBind(
					queue.Name,                   // queue name
					routingKey,                   // routing key
					c.RabbitMQ.Producer.Exchange, // exchange
					false,
					nil)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	var queue amqp.Queue
	if c.RabbitMQ.Consumer != nil {
		if c.RabbitMQ.Consumer.Queue != "" {
			queue, err = channel.QueueDeclare(
				c.RabbitMQ.Consumer.Queue, // name
				true,                      // durable
				false,                     // delete when unused
				false,                     // exclusive
				false,                     // no-wait
				nil,                       // arguments
			)
			if err != nil {
				return nil, err
			}

		}
	}

	return &RabbitMQ{
		Connection:      conn,
		Channel:         channel,
		ProdExchange:    exchange,
		ProdRoutingKeys: prodRks,
		Queue:           &queue,
	}, nil
}

func (r *RabbitMQ) Close() error {
	if err := r.Connection.Close(); err != nil {
		return err
	}
	if err := r.Channel.Close(); err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) StartConsume(queueName string, handlers ...util.RabbitMQMsgHandleFunc) error {
	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	// tạo channel để chờ vô hạn
	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			for _, handler := range handlers {
				handler(msg)
			}
		}
	}()

	<-forever

	return nil
}

func (r *RabbitMQ) Publish(exchange string, routingKey string, header amqp.Table, body interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return r.Channel.Publish(
		exchange,                                   // exchange
		fmt.Sprintf("%v-%v", exchange, routingKey), // routing key
		false,                                      // mandatory
		false,                                      // immediate
		amqp.Publishing{
			DeliveryMode:    amqp.Persistent,
			Headers:         header,
			ContentType:     "application/json",
			ContentEncoding: constant.CharSetUtf8,
			Body:            jsonBody,
		})
}

func (r *RabbitMQ) PublishWithRouting(exchange string, routing RabbitMQRouting, header amqp.Table, body interface{}) error {
	return r.Publish(exchange, routing.GetRoutingKey(r.ProdRoutingKeys, body), header, body)
}

type RabbitMQRouting interface {
	GetRoutingKey(routingKeys map[string]string, data interface{}) string
}
