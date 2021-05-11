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
		text: `    Lorem ipsum dolor sit amet consectetur adipisicing elit. Veniam consequatur fugit eligendi sunt voluptate vel perspiciatis beatae numquam dignissimos, commodi, cupiditate unde, repudiandae doloribus nisi minima rem placeat modi nesciunt.
		Illo porro dignissimos cum quidem magnam ab obcaecati quibusdam debitis labore ad iusto atque, reprehenderit autem pariatur nam ipsam possimus adipisci quae sunt architecto aperiam totam unde hic. Perspiciatis, quia.
		Laboriosam laudantium dolore iste quam adipisci molestias quae. Laborum impedit consectetur maiores velit nihil aperiam inventore fuga amet, quasi, quidem, recusandae aut saepe odit voluptas rerum? Assumenda odit ullam unde?
		Ducimus ipsam, aspernatur eveniet eius eum cumque molestiae ab consequuntur odit temporibus cum, adipisci beatae nobis in corporis eligendi commodi aperiam et officiis est dolor, totam aliquid! Quisquam, beatae est.
		Delectus expedita accusamus ducimus consequatur quae quibusdam cum placeat illum, id explicabo deserunt repudiandae fugit consequuntur voluptatem porro asperiores amet voluptates rem cupiditate. Ipsa, sed. Pariatur harum molestias impedit corrupti?`,
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
