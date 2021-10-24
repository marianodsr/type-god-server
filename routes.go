package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleRoutes() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from server")
	})
	r.Get("/ws", serveWs)

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connected!")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "unable to connect", http.StatusInternalServerError)
		return
	}

	go S.newClient(ws)

}
