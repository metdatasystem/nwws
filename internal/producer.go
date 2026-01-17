package nwws

import (
	"context"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type Message struct {
	contentType string
	data        []byte
}

type Producer struct {
	channel  *amqp.Channel
	messages chan Message
	done     bool
}

func NewProducer() (*Producer, error) {

	conn, err := amqp.Dial(os.Getenv("RABBIT_URL"))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"nws.queue", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	producer := &Producer{
		channel:  ch,
		messages: make(chan Message),
		done:     false,
	}
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func (p *Producer) Run() {
	for message := range p.messages {
		if p.done {
			p.channel.Close()
			return
		}

		err := p.SendMessage(message)
		if err != nil {
			log.Error().Err(err).Msg("failed to send message")
		}

	}
}

func (p *Producer) Stop() {
	p.done = true
}

func (p *Producer) SendMessage(message Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("message")

	return p.channel.PublishWithContext(ctx,
		"",          // exchange
		"nws.queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: message.contentType,
			Timestamp:   time.Now(),
			AppId:       "us.ingest.awips",
			Body:        message.data,
		})
}
