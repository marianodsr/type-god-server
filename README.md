# type-god-server

Server for type-god.

Opens a websocket connection for every client.

Creates rooms and when filled, it starts the game notifiying the clients and serving them the text piece to type.

By achieving concurrency with goroutines and channel communication, it can handle large numbers of connections with ease.
