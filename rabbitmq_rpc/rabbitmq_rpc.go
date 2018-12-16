package rabbitmq_rpc

import (
	"os"

	"github.com/streadway/amqp"
	"github.com/taeduard/PI-SMS/utils"
)

//"amqp://guest:guest@localhost:5672/"
func CheckAMQP_URL() (amqp string) {
	amqp = os.Getenv("AMQP_URL")
	if amqp != "" {
		return amqp
	} else {
		return "amqp://guest:guest@localhost:5672/"
	}
}

func RPCcommand(command string) (res string, err error) {
	conn, err := amqp.Dial(CheckAMQP_URL())
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")
	corrId := utils.RandomString(32)
	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(command),
		})
	utils.FailOnError(err, "Failed to publish a message")
	for d := range msgs {
		if corrId == d.CorrelationId {
			res = string(d.Body)
			break
		}
	}
	return
}
