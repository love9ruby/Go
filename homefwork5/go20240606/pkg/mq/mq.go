package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

var URL = "amqp://guest:guest@localhost:5672/"

type IMQ interface {
	// Publish publish a message to a queue
	Publish(message []byte) error
	// Consume consume a message from a queue
	Consume() (<-chan amqp.Delivery, error)
	// Close destroy the connection
	Close() error
}

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	URL       string
}

func (r *RabbitMQ) Publish(message []byte) error {
	err := r.channel.Publish(r.Exchange, r.QueueName, false, false, amqp.Publishing{
		ContentType: "json/application",
		Body:        message,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	msgChan, err := r.channel.Consume(r.QueueName, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return msgChan, nil
}

func NewIMQ() (rMQ IMQ, err error) {
	host := os.Getenv("RABBITMQ_HOST")
	if host != "" {
		URL = fmt.Sprintf("amqp://guest:guest@%s:5672/", host)
	}
	rMQ, err = NewRabbitMQ("userEvent", "", "")
	if err != nil {
		return nil, err
	}
	return rMQ, nil
}

func NewRabbitMQ(queueName string, exchange string, key string) (*RabbitMQ, error) {
	rmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, URL: URL}
	var err error
	rmq.conn, err = amqp.Dial(rmq.URL)
	if err != nil {
		return nil, err
	}
	rmq.channel, err = rmq.conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = rmq.channel.QueueDeclare(rmq.QueueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return rmq, nil
}

func (r *RabbitMQ) Close() error {
	err := r.channel.Close()
	if err != nil {
		return err
	}
	err = r.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
