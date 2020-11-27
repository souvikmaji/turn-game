package main

import (
	"log"
)

// Room manages the set of active clients and broadcasts messages to the
// clients.
type Room struct {
	// Register requests from the clients.
	register chan *Client

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte
}

func newRoom() *Room {
	return &Room{
		register:  make(chan *Client),
		broadcast: make(chan []byte),
		clients:   make(map[*Client]bool),
	}
}

func (r *Room) run() {
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
