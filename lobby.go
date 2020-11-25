package main

import (
	"log"
)

// Lobby maintains the set of active clients and broadcasts messages to the
// clients.
type Lobby struct {
	// Register requests from the clients.
	register chan *Client

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte
}

func newLobby() *Lobby {
	return &Lobby{
		register:  make(chan *Client),
		broadcast: make(chan []byte),
		clients:   make(map[*Client]bool),
	}
}

func (r *Lobby) run() {
	for {
		select {
		case client := <-r.register:
			// adding new client to the registered clients map
			r.clients[client] = true
		case message := <-r.broadcast:
			log.Printf("broadcasting message to clients: %s\n", string(message))
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}
