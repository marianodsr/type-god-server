package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type client struct {
	conn *websocket.Conn
	nick string
	room *room
	data chan []byte
}

func (c *client) readMessages() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			c.conn.Close()
			break
		}
		c.room.forward <- msg
	}

}

func (c *client) sendMsg(msg string) {
	if err := c.conn.WriteMessage(1, []byte(msg)); err != nil {
		fmt.Printf("unable to write message, err: %s", err.Error())
		c.conn.Close()

	}
}
