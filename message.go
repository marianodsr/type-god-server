package main

type Event string

const (
	JOIN_ROOM    Event = "JOIN_ROOM"
	START_GAME   Event = "START_GAME"
	NOTIFICATION Event = "NOTIFICATION"
)

type Message struct {
	Event Event                  `json:"event"`
	Data  map[string]interface{} `json:"data"`
}
