package rabbitmq

import (
	"koffee/pkg/pool"

	"github.com/streadway/amqp"
)

// this simple file and interfaces work like this:
// ** For initializing
// r := rabbitmq.Initialize()
// defer r.Stop()
// ** For listening into a queue in another thread
// r.OnMessage("test", func(m *amqp.Delivery) {
// 	fmt.Println("New message!", string(m.Body))
// })
// ** For initializing a queue and getting a sender to that queue
// s := r.NewSender("test")
// ** For sending a message to that queue
// s.Send("Hello!", "text/plain")

// ** It's organized like this so it's easier to pass this dependency into controllers

// messageListenerRabbitMQ is for listening messages and can produce senders for RabbitMQ implementing MessageListener
type messageListenerRabbitMQ struct {
	channel *amqp.Channel
	pool    *pool.Pool
}

// messageSenderRabbitMQ implements MessageSender for RabbitMQ
type messageSenderRabbitMQ struct {
	queueName string
	channel   *amqp.Channel
}

// MessageSender interface is for sending messages to the desired queue or channel
type MessageSender interface {
	Send(b []byte, contentType string)
}

// MessageListener is for listening messages and can produce senders
type MessageListener interface {
	NewSender(name string) MessageSender
	OnMessage(queueName string, action func(msg *amqp.Delivery))
	Stop()
}

// Send sends a message through MessageSender
func (m *messageSenderRabbitMQ) Send(b []byte, contentType string) {
	m.channel.Publish("", m.queueName, false, false, amqp.Publishing{ContentType: contentType, Body: b})
}

// NewSender Creates a new sender for sending messages to the specified queue
func (m *messageListenerRabbitMQ) NewSender(name string) MessageSender {
	q, err := m.channel.QueueDeclare(name, false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}
	return &messageSenderRabbitMQ{queueName: q.Name, channel: m.channel}
}

// OnMessage executes an action when there is a message coming from the queue
func (m *messageListenerRabbitMQ) OnMessage(queueName string, action func(msg *amqp.Delivery)) {
	queue, err := m.channel.QueueDeclare(queueName, false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}
	msgs, err := m.channel.Consume(queue.Name,
		"",
		// consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		panic(err)
	}
	m.pool.Schedule(func() {
		for msg := range msgs {
			// Capture value
			msg := msg
			m.pool.Schedule(func() { action(&msg) })
		}
	})
}

func (m *messageListenerRabbitMQ) Stop() {
	m.Stop()
}

// Initialize inits RabbitMQ subscription
func Initialize() MessageListener {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil
	}
	channel, err := conn.Channel()
	rabbitPool := pool.NewPool(50, 25, 25)
	return &messageListenerRabbitMQ{pool: rabbitPool, channel: channel}
}
