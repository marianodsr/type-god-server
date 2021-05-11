package main

import (
	"github.com/google/uuid"
)

type room struct {
	id      string
	members map[*client]bool
	join    chan *client
	leave   chan *client
	forward chan []byte
}

func newRoom() *room {
	r := &room{
		id:      uuid.New().String(),
		members: make(map[*client]bool),
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
	}
	go r.run()
	return r
}

func (r *room) broadcast(msg []byte) {
	for client := range r.members {
		client.sendMsg(string(msg))
	}
}

func (r *room) run() {
	for {
		select {
		case msg := <-r.forward:
			r.broadcast(msg)
		case client := <-r.join:
			r.members[client] = true
		}
	}
}
