package main

import (
	"github.com/google/uuid"
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
		conn:     ws,
		data:     make(chan []byte),
		nick:     "guest " + uuid.New().String(),
		progress: 0,
	}
	s.addToRoom(c)
	c.readMessages()
}

func (s *server) addToRoom(c *client) {

	if len(s.rooms) == 0 {
		r := newRoom(3)
		r.join <- c
		s.rooms = append(s.rooms, *r)
		return
	}
	for _, room := range s.rooms {
		if len(room.members) < 3 {
			room.join <- c
			return
		}
	}
	r := newRoom(3)
	r.join <- c
	s.rooms = append(s.rooms, *r)

}
