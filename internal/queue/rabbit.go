package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mohdhujaifa/profile/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

const ContactExchange = "portfolio.events"
const ContactRoutingKey = "contact.created"

type Publisher interface {
	PublishContact(ctx context.Context, msg model.ContactMessage) error
	Close() error
}

type RabbitPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbit(url string) (*RabbitPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}
	if err := ch.ExchangeDeclare(ContactExchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq exchange: %w", err)
	}
	return &RabbitPublisher{conn: conn, channel: ch}, nil
}

func (p *RabbitPublisher) PublishContact(ctx context.Context, msg model.ContactMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.channel.PublishWithContext(ctx, ContactExchange, ContactRoutingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func (p *RabbitPublisher) Close() error {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

type NoopPublisher struct{}

func (NoopPublisher) PublishContact(_ context.Context, _ model.ContactMessage) error {
	return nil
}

func (NoopPublisher) Close() error { return nil }
