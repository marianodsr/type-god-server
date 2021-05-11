package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type client struct {
	conn     *websocket.Conn
	nick     string
	room     *room
	data     chan []byte
	progress int
}

func (c *client) readMessages() {
	for {
		msg := Message{}
		err := c.conn.ReadJSON(msg)
		if err != nil {
			c.conn.Close()
			break
		}
		c.room.forward <- msg
	}

}

func (c *client) sendMsg(msg Message) {
	if err := c.conn.WriteJSON(msg); err != nil {
		fmt.Printf("unable to write message, err: %s", err.Error())
		c.conn.Close()

	}
}
