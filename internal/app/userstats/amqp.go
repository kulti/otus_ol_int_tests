package userstats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type amqpConn struct {
	conn   *amqp.Connection
	amqpCh *amqp.Channel
	msgs   <-chan amqp.Delivery
	store  *store
}

func connectAmqp(store *store) (*amqpConn, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		"chess.game.stats", // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exhange: %w", err)
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(
		q.Name,             // queue name
		"finish",           // routing key
		"chess.game.stats", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return &amqpConn{
		conn:   conn,
		amqpCh: ch,
		msgs:   msgs,
		store:  store,
	}, nil

}

func (c *amqpConn) Run() {
	var msg struct {
		GameID  string `json:"game_id"`
		UserID  string `json:"user_id"`
		Outcome string `json:"outcome"`
	}

	for d := range c.msgs {
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			log.Println("failed to unmarshal msg: ", err)
			continue
		}

		err := c.store.SaveUserStats(context.Background(), msg.UserID, msg.GameID, msg.Outcome)
		if err != nil {
			log.Println("failed to save msg: ", err)
			continue
		}

		log.Printf("save msg %+v\n", msg)
	}
}

func (c *amqpConn) Close() {
	c.amqpCh.Close()
	c.conn.Close()
}
