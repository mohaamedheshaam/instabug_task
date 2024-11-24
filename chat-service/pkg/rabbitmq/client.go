package rabbitmq

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/streadway/amqp"
)

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
}

type Client struct {
    connection *amqp.Connection
    channel    *amqp.Channel
}

func NewClient(cfg Config) (*Client, error) {
    dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", 
        cfg.User, cfg.Password, cfg.Host, cfg.Port)

    conn, err := amqp.Dial(dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to create RabbitMQ channel: %w", err)
    }

    queues := []string{"chat_created", "message_created"}
    for _, queue := range queues {
        if err := ch.ExchangeDeclare(
            queue,   // name
            "topic", // type
            true,    // durable
            false,   // auto-deleted
            false,   // internal
            false,   // no-wait
            nil,     // arguments
        ); err != nil {
            conn.Close()
            return nil, fmt.Errorf("failed to declare exchange %s: %w", queue, err)
        }
    }

    return &Client{
        connection: conn,
        channel:    ch,
    }, nil
}

func (c *Client) Close() error {
    if err := c.channel.Close(); err != nil {
        return fmt.Errorf("failed to close channel: %w", err)
    }
    if err := c.connection.Close(); err != nil {
        return fmt.Errorf("failed to close connection: %w", err)
    }
    return nil
}

func (c *Client) PublishChatCreated(ctx context.Context, data interface{}) error {
    body, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal chat data: %w", err)
    }

    return c.publish("chat_created", body)
}

func (c *Client) PublishMessageCreated(ctx context.Context, data interface{}) error {
    body, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal message data: %w", err)
    }

    return c.publish("message_created", body)
}

func (c *Client) publish(queue string, body []byte) error {
    return c.channel.Publish(
        queue, // exchange
        "",    // routing key
        false, // mandatory
        false, // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:        body,
            DeliveryMode: amqp.Persistent,
            Timestamp:   time.Now(),
        },
    )
}
