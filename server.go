package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type server struct {
	rooms []room
}

func newServer() *server {
	return &server{}
}

func (s *server) newClient(ws *websocket.Conn) {
	c := &client{
		conn: ws,
		data: make(chan []byte),
		nick: "guest",
	}
	s.addToRoom(c)
	c.readMessages()
}

func (s *server) addToRoom(c *client) {
	if len(s.rooms) == 0 {
		r := newRoom()
		r.join <- c
		c.room = r
		s.rooms = append(s.rooms, *r)
		c.sendMsg(fmt.Sprintf("%s joined to room %s", c.nick, r.id))
		return
	}
	for _, room := range s.rooms {
		if len(room.members) < 3 {
			room.join <- c
			c.room = &room
			c.sendMsg(fmt.Sprintf("%s joined to room %s", c.nick, room.id))

			return
		}
	}

	r := newRoom()
	r.join <- c
	c.room = r
	c.sendMsg(fmt.Sprintf("%s joined to room %s", c.nick, r.id))
	s.rooms = append(s.rooms, *r)

}
