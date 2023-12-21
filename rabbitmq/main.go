package rabbitmq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
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
