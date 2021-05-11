package main

var S *server

func main() {
	S = newServer()
	handleRoutes()
}
