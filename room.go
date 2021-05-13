package main

import (
	"fmt"

	"github.com/google/uuid"
)

type room struct {
	id       string
	members  map[*client]bool
	join     chan *client
	leave    chan *client
	forward  chan Message
	text     string
	capacity int
}

func newRoom(capacity int) *room {
	r := &room{
		id:       uuid.New().String(),
		members:  make(map[*client]bool),
		forward:  make(chan Message),
		join:     make(chan *client),
		leave:    make(chan *client),
		capacity: capacity,
		text:     `At that moment he had a thought that he'd never imagine he'd consider. "I could just cheat," he thought, "and that would solve the problem." He tried to move on from the thought but it was persistent. It didn't want to go away and, if he was honest with himself, he didn't want it to.`,
	}
	go r.run()
	return r
}

func (r *room) broadcast(msg Message) {
	for client := range r.members {
		client.sendMsg(msg)
	}
}

func (r *room) run() {
	for {
		select {
		case msg := <-r.forward:
			r.broadcast(msg)
		case client := <-r.join:
			fmt.Printf("\n%s client s has joined the room", client.nick)
			r.members[client] = true
			client.room = r
			msg := Message{
				Event: NOTIFICATION,
				Data: map[string]interface{}{
					"msg": fmt.Sprintf("%s joined room %s", client.nick, r.id),
				},
			}
			client.sendMsg(msg)
			if len(r.members) == r.capacity {
				go r.startGame()
			}
		case client := <-r.leave:
			fmt.Printf("\nclient %s has left the room", client.nick)
			delete(r.members, client)
			client.conn.Close()
		}
	}
}

func (r *room) startGame() {
	var players []string
	for player := range r.members {
		players = append(players, player.nick)
	}
	msg := Message{
		Event: START_GAME,
		Data: map[string]interface{}{
			"text":    r.text,
			"players": players,
		},
	}
	r.forward <- msg
	fmt.Printf("Game started on room %s", r.id)

}
