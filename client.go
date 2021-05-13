package main

import (
	"fmt"
	"math"

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
		err := c.conn.ReadJSON(&msg)
		fmt.Printf("\n%+v", msg)
		if err != nil {
			c.room.leave <- c
			c.conn.Close()
			break
		}
		switch msg.Event {
		case PROGRESS_ADVANCE:
			floatProgress := (float64(len(msg.Data["doneText"].(string))) / float64(len(c.room.text))) * 100
			intProgress := math.Floor(floatProgress)
			fmt.Println(intProgress)
			fmt.Println(floatProgress)
			msg.Data = map[string]interface{}{
				"progress": intProgress,
				"player":   c.nick,
			}
			c.room.forward <- msg
		}
	}

}

func (c *client) sendMsg(msg Message) {
	if err := c.conn.WriteJSON(msg); err != nil {
		fmt.Printf("unable to write message, err: %s", err.Error())
		c.conn.Close()

	}
}
