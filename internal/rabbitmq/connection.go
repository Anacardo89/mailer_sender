package rabbitmq

import (
	"github.com/Anacardo89/mailer_sender/pkg/logger"
	"github.com/Anacardo89/mailer_sender/pkg/rabbit"
	"github.com/streadway/amqp"
)

func StartWorkers(r *rabbit.Config, conn *amqp.Connection, ch *amqp.Channel, msgs chan<- amqp.Delivery) {
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
