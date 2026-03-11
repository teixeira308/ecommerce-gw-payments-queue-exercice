package broker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	var conn *amqp.Connection
	var err error

	// Retry loop for RabbitMQ connection
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/10): %v. Retrying in 5s...", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQClient{conn: conn, ch: ch}, nil
}

func (c *RabbitMQClient) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *RabbitMQClient) Publish(ctx context.Context, exchange, routingKey string, body interface{}) error {
	log.Printf("[RabbitMQ] Publishing message to exchange '%s' with routing key '%s'", exchange, routingKey)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("[RabbitMQ] Error marshaling message for exchange '%s': %v", exchange, err)
		return err
	}

	err = c.ch.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonBody,
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		log.Printf("[RabbitMQ] Error publishing message to exchange '%s': %v", exchange, err)
		return err
	}

	return nil
}

func (c *RabbitMQClient) Consume(queueName, consumerName string, handler func(d amqp.Delivery)) error {
	log.Printf("[RabbitMQ] Starting consumer '%s' for queue '%s'", consumerName, queueName)
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
		log.Printf("[RabbitMQ] Error starting consumer '%s' for queue '%s': %v", consumerName, queueName, err)
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("[RabbitMQ] Message received by consumer '%s' from queue '%s'", consumerName, queueName)
			handler(d)
		}
		log.Printf("[RabbitMQ] Consumer '%s' stopped for queue '%s'", consumerName, queueName)
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
