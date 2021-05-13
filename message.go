package main

type Event string

const (
	JOIN_ROOM        Event = "JOIN_ROOM"
	START_GAME       Event = "START_GAME"
	NOTIFICATION     Event = "NOTIFICATION"
	PROGRESS_ADVANCE Event = "PROGRESS_ADVANCE"
)

type Message struct {
	Event Event                  `json:"event"`
	Data  map[string]interface{} `json:"data"`
}
