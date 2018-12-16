package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/taeduard/PI-SMS/rabbitmq_rpc"
	"github.com/taeduard/PI-SMS/utils"
)

func main() {
	conn, err := amqp.Dial(rabbitmq_rpc.CheckAMQP_URL())
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		i := 0
		for d := range msgs {
			i++
			time := time.Now().UTC()
			response := utils.Execute(string(d.Body))
			hostname := utils.Execute("hostname")
			// fmt.Printf("%s\n:\n%s:\n%s", time.String(), hostname, response)
			err := ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(fmt.Sprintf("%s\n;\n%s;\n%s", time.String(), hostname, response)),
				})
			utils.FailOnError(err, "Failed to publish a message")
			d.Ack(false)
		}
	}()
	<-forever
}
