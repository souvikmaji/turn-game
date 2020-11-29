package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

var (
	maxRoomSize = 4
)

// Room manages the set of active clients and broadcasts messages to the
// clients.
type Room struct {
	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Registered clients and their scores.
	clients map[*Client]int

	// keeping track if more players can join this room
	isFull bool

	// Inbound messages from the clients
	broadcast chan *Client
}

func newRoom() *Room {
	r := &Room{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Client),
		clients:    make(map[*Client]int),
		isFull:     false,
	}

	go r.run()
	return r
}

func (r *Room) run() {
	rand.Seed(time.Now().UnixNano())

	for {
		select {
		case client := <-r.register:
			// adding new client to the registered clients map
			r.clients[client] = 0
			roomSize := len(r.clients)

			client.position = roomSize
			if roomSize == maxRoomSize {
				r.isFull = true
			}
		case client := <-r.unregister:
			// removing client once disconnected
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case client := <-r.broadcast:

			score := rand.Intn(6)
			r.clients[client] = score

			message, err := json.Marshal(getScores(r.clients))
			if err != nil {
				log.Println("error marshalling scores", err)
				// TODO: send error message to client
				break
			}

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

func getScores(clients map[*Client]int) map[string]int {
	scores := make(map[string]int)

	for client, score := range clients {
		scores[client.name] = score
	}
	return scores
}
