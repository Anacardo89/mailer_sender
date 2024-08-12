package main

import (
	"log"

	"github.com/Anacardo89/mailer_sender/internal/config"
	"github.com/Anacardo89/mailer_sender/internal/handlers"
	"github.com/Anacardo89/mailer_sender/internal/logger"
	"github.com/streadway/amqp"
)

func main() {
	logger.CreateLogger()

	// Rabbit Setup
	rabbit := config.LoadRabbitConfig()

	conn := rabbit.Connect()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error.Fatal(err)
	}
	defer ch.Close()

	rabbit.DeclareQueues(ch)

	msgs := make(chan amqp.Delivery)
	rabbit.StartWorkers(conn, ch, msgs)

	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			switch msg.RoutingKey {
			case "register_email":
				handlers.SendRegistrationEmail(msg)
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
