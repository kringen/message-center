package rabbitmq

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

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
		logger.Info(fmt.Sprintf("Waiting for RabbitMQ Connection.  Attempt %d\n", i+1))
		if conn == nil {
			time.Sleep(interval)
		}
	}
	if conn != nil {
		logger.Info("Connection is successful!")
		m.Connection = conn
		m.Channel, err = m.Connection.Channel()
		if err != nil {
			panic(err)
		}
		return nil
	} else {
		logger.Info("Connection timed out")
		return err
	}

}

func (m *MessageCenter) CreateQueue(name string, durable bool, deleteUnused bool,
	exclusive bool, noWait bool, arguments map[string]interface{}) error {
	logger.Info(fmt.Sprintf("Creating queue: %s", name))
	_, err := m.Channel.QueueDeclare(name, durable, deleteUnused, exclusive, noWait, arguments)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageCenter) ReceiveMessage(replyChannel chan string, queue string) {
	// Listen for messages
	messages, err := m.ConsumeMessage(queue, "", true, false, false, false, nil)
	if err != nil {
		logger.Error(err.Error())
	}
	for message := range messages {
		// Continue to receive messages
		logger.Info(fmt.Sprintf("Received message on %s", queue))
		logger.Debug(fmt.Sprintf("Received message on %s: %s", queue, message.Body))
		replyChannel <- string(message.Body)
	}
}

func (m *MessageCenter) PublishMessage(queue string, message []byte, exchange string, mandatory bool, immediate bool, contentType string, correlationId string, replyTo string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger.Info(fmt.Sprintf("Publishing message on queue: %s", queue))
	logger.Debug(fmt.Sprintf("Publishing message %v on queue: %s", message, queue))
	err := m.Channel.PublishWithContext(ctx,
		exchange,  // exchange ""
		queue,     // routing key
		mandatory, // mandatory
		immediate, // immediate
		amqp.Publishing{
			CorrelationId: correlationId,
			ReplyTo:       replyTo,
			ContentType:   contentType, //"text/plain"
			Body:          message,
		})

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (m *MessageCenter) ConsumeMessage(queue string, consumer string, autoAck bool, exclusive bool,
	noLocal bool, noWait bool, arguments map[string]interface{}) (<-chan amqp.Delivery, error) {
	logger.Info(fmt.Sprintf("Consuming on queue: %s", queue))
	messages, err := m.Channel.Consume(
		queue,     // queue name
		consumer,  // consumer ""
		autoAck,   // auto-ack true
		exclusive, // exclusive false
		noLocal,   // no local false
		noWait,    // no wait false
		arguments, // arguments nil
	)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

