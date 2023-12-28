package rabbitmq

import (
	"log"
	"time"

	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageCenter struct {
	ServerUrl  string
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (m *MessageCenter) Connect(channel string, attempts int, intervalSeconds int) error {

	interval := time.Duration(intervalSeconds) * time.Second
	var conn *amqp.Connection
	var err error
	for i := 0; i < attempts; i++ {
		conn, err = amqp.Dial(m.ServerUrl)
		log.Printf("Waiting for RabbitMQ Connection.  Attempt %d\n", i+1)
		if conn == nil {
			time.Sleep(interval)
		}
	}
	if conn != nil {
		log.Printf("Connection is successful!")
		m.Connection = conn
		m.Channel, err = m.Connection.Channel()
		if err != nil {
			panic(err)
		}
		return nil
	} else {
		log.Printf("Connection Timed Out!")
		return err
	}

}

func (m *MessageCenter) CreateQueue(name string, durable bool, deleteUnused bool,
	exclusive bool, noWait bool, arguments map[string]interface{}) error {

	_, err := m.Channel.QueueDeclare(name, durable, deleteUnused, exclusive, noWait, arguments)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageCenter) PublishMessage(ctx context.Context, queue string, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var err error
	err = m.Channel.PublishWithContext(ctx,
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		return err
	} else {
		return nil
	}
}
