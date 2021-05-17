package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/google/uuid"
)

type room struct {
	mu             *sync.Mutex
	id             string
	members        map[*client]bool
	join           chan *client
	leave          chan *client
	forward        chan Message
	text           string
	capacity       int
	gameInProgress bool
}

type textPiece struct {
	Paragraph string `json:"paragraph"`
}

func newRoom(capacity int) *room {

	f, err := os.Open("paragraphs.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	var textCollection []textPiece
	if err := json.NewDecoder(f).Decode(&textCollection); err != nil {
		log.Fatalln(err)
	}

	text := textCollection[rand.Intn(len(textCollection))]

	r := &room{
		id:             uuid.New().String(),
		members:        make(map[*client]bool),
		forward:        make(chan Message),
		join:           make(chan *client),
		leave:          make(chan *client),
		capacity:       capacity,
		mu:             new(sync.Mutex),
		gameInProgress: false,
		text:           text.Paragraph,
	}
	go r.run() // starts listening on room channels0

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
			r.mu.Lock()
			r.members[client] = true
			r.mu.Unlock()
			client.room = r
			msg := Message{
				Event: NOTIFICATION,
				Data: map[string]interface{}{
					"msg": fmt.Sprintf("%s joined room %s", client.nick, r.id),
				},
			}
			client.sendMsg(msg)
			if len(r.members) == r.capacity {
				go r.startGame() // needs to be concurrent in order to send message to forward channel
			}
		case client := <-r.leave:
			fmt.Printf("\nclient %s has left the room", client.nick)
			r.mu.Lock()
			delete(r.members, client)
			r.mu.Unlock()
			client.conn.Close()
		}
	}
}

func (r *room) startGame() {
	var players []string
	r.mu.Lock()
	for player := range r.members {
		players = append(players, player.nick)
	}
	r.mu.Unlock()
	msg := Message{
		Event: START_GAME,
		Data: map[string]interface{}{
			"text":    r.text,
			"players": players,
		},
	}
	r.mu.Lock()
	r.gameInProgress = true
	r.mu.Unlock()
	r.forward <- msg
	fmt.Printf("Game started on room: %s", r.id)

}
