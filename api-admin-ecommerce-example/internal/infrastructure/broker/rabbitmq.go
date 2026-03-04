package broker

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{conn: conn, ch: ch}, nil
}

func (c *RabbitMQClient) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *RabbitMQClient) Publish(ctx context.Context, exchange, routingKey string, body interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return c.ch.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonBody,
			DeliveryMode: amqp.Persistent,
		})
}

func (c *RabbitMQClient) Consume(queueName, consumerName string, handler func(d amqp.Delivery)) error {
	msgs, err := c.ch.Consume(
		queueName,    // queue
		consumerName, // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
		log.Println("RabbitMQ consumer stopped")
	}()

	return nil
}

func (c *RabbitMQClient) SetupTopology() error {
	// Exchange
	err := c.ch.ExchangeDeclare(
		"payments.exchange", // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return err
	}

	// Queues
	_, err = c.ch.QueueDeclare(
		"payment.requested.queue", // name
		true,                      // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		return err
	}

	_, err = c.ch.QueueDeclare(
		"payment.processed.queue", // name
		true,                      // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		return err
	}

	_, err = c.ch.QueueDeclare(
		"payments.dlq", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    "payments.exchange",
			"x-dead-letter-routing-key": "payment.dead",
		},
	)
	if err != nil {
		return err
	}

	// Bindings
	err = c.ch.QueueBind(
		"payment.requested.queue", // queue name
		"payment.requested",       // routing key
		"payments.exchange",       // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = c.ch.QueueBind(
		"payment.processed.queue", // queue name
		"payment.processed",       // routing key
		"payments.exchange",       // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = c.ch.QueueBind(
		"payments.dlq",      // queue name
		"payment.dead",      // routing key
		"payments.exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
