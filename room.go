package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

var (
	maxRoomSize = 4
	scoreToWin  = 10 // TODO: change at last based on requirement
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

	nextTurn int

	// keeping track if more players can join this room
	isFull bool

	// current winner of this room
	winner *Client

	// Inbound messages from the clients
	broadcast chan *Client
}

func newRoom() *Room {
	r := &Room{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Client),
		clients:    make(map[*Client]int),
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

			client.position = roomSize - 1

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
			if len(r.clients) < 2 {
				log.Println("waiting for other players to join. ignoring message")
				break
			}
			r.nextTurn++
			r.nextTurn = r.nextTurn % len(r.clients)

			// generate new client score
			score := rand.Intn(7)
			r.clients[client] += score

			if r.clients[client] >= scoreToWin {
				r.winner = client
			}

			// create client response
			message, err := r.createResponse()
			if err != nil {
				log.Println("error marshalling scores", err)
				// TODO: send error message to client
				break
			}

			// broadcast response to all clients in room
			log.Printf("broadcasting message to clients: %s\n", string(message))

			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}

			// reset room once winner is declared
			if r.winner != nil {
				for client := range r.clients {
					r.clients[client] = 0
				}
			}
			r.winner = nil
		}
	}
}

// setup next turn
func (r *Room) setNextTurn() {
	r.nextTurn++
	r.nextTurn = r.nextTurn % len(r.clients)
}

func (r *Room) createResponse() ([]byte, error) {
	response := newResponse()

	if r.winner != nil {
		response.Winner = r.winner.username
	}

	response.setScores(r.clients)

	return json.Marshal(response)
}
