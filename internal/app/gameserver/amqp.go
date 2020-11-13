package gameserver

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type amqpConn struct {
	conn   *amqp.Connection
	amqpCh *amqp.Channel
}

func connectAmqp() (*amqpConn, error) {
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

	return &amqpConn{
		conn:   conn,
		amqpCh: ch,
	}, nil

}

type msg struct {
	GameID  string `json:"game_id"`
	UserID  string `json:"user_id"`
	Outcome string `json:"outcome"`
}

func (c *amqpConn) Send(m msg) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return c.amqpCh.Publish(
		"chess.game.stats", //exchange
		"finish",           //key
		false,              //mandatory
		false,              //immediate
		amqp.Publishing{
			Body: b,
		},
	)
}

func (c *amqpConn) Close() {
	c.amqpCh.Close()
	c.conn.Close()
}
