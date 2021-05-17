package main

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type server struct {
	mu    *sync.Mutex
	rooms []*room
}

func newServer() *server {

	return &server{
		mu: new(sync.Mutex),
	}
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

// func (s *server) addToRoom(c *client) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if len(s.rooms) == 0 {
// 		r := newRoom(3)
// 		r.join <- c
// 		s.rooms = append(s.rooms, r)

// 		return
// 	}
// 	for _, room := range s.rooms {
// 		room.mu.Lock()
// 		if len(room.members) < 3 && !room.gameInProgress {
// 			room.join <- c
// 			room.mu.Unlock()
// 			return
// 		}
// 		room.mu.Unlock()
// 	}

// }

func (s *server) addToRoom(c *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.rooms) == 0 {
		r := newRoom(3)
		r.join <- c
		s.rooms = append(s.rooms, r)

		return
	}
	for _, room := range s.rooms {
		room.mu.Lock()
		if len(room.members) < 3 && !room.gameInProgress {
			room.join <- c
			room.mu.Unlock()
			return
		}
		room.mu.Unlock()
	}
	r := newRoom(3)
	r.join <- c
	s.rooms = append(s.rooms, r) //CHECK IF THREAD SAFE

}
