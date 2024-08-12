package rabbitmq

import (
	"fmt"

	"github.com/Anacardo89/mailer_sender/internal/logger"
	"github.com/streadway/amqp"
)

type Config struct {
	RabbitUser string   `yaml:"rabbit_user"`
	RabbitPass string   `yaml:"rabbit_pass"`
	RabbitHost string   `yaml:"rabbit_host"`
	RabbitPort string   `yaml:"rabbit_port"`
	Queues     []string `yaml:"queues"`
}

func (r *Config) Connect() *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s%s/",
		r.RabbitUser, r.RabbitPass, r.RabbitHost, r.RabbitPort)
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error.Fatal(err)
	}
	return conn
}

func (r *Config) DeclareQueues(ch *amqp.Channel) {
	for _, queue := range r.Queues {
		_, err := ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			logger.Error.Fatal(err)
		}
	}
}

func (r *Config) StartWorkers(conn *amqp.Connection, ch *amqp.Channel, msgs chan<- amqp.Delivery) {
	for _, queue := range r.Queues {
		go func(q string) {
			worker(ch, q, msgs)
		}(queue)
	}
}

func worker(ch *amqp.Channel, queue string, msgs chan<- amqp.Delivery) {
	msg, err := ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logger.Error.Fatal(err)
	}
	for m := range msg {
		msgs <- m
	}
}
