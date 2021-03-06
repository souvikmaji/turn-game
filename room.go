package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"
)

var (
	maxRoomSize = 4
	scoreToWin  = 61
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
	playerMove chan *Client
}

func newRoom() *Room {
	r := &Room{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		playerMove: make(chan *Client),
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
				// if it was the disconnected players turn, pass the turn to the next player
				if r.nextTurn == client.position {
					r.nextTurn = r.getNextTurn()
				}
				delete(r.clients, client)
				close(client.send)
			}
		case client := <-r.playerMove:

			// min two players required to play game
			if len(r.clients) < 2 {
				r.handleError(client, errors.New("waiting for other players to join"))
				break
			}

			//  allow only when next turn is for this client
			if client.position != r.nextTurn {
				r.handleError(client, errors.New("Not your turn"))
				break
			}

			r.nextTurn = r.getNextTurn()

			// generate new client score
			score := rand.Intn(7)
			r.clients[client] += score

			if r.clients[client] >= scoreToWin {
				r.winner = client
			}

			// create client response
			message, err := r.createResponse()
			if err != nil {
				r.handleError(client, err)
				break
			}

			// broadcast response to all clients in room
			log.Printf("broadcasting message to clients: %s\n", string(message))

			r.broadcast(message)

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
func (r *Room) getNextTurn() int {
	return (r.nextTurn + 1) % len(r.clients)
}

// check if next turn is valid
func (r *Room) isValidNextTurn(nextMove int) bool {
	for client := range r.clients {
		if client.position == nextMove {
			return true
		}
	}
	return false
}

func (r *Room) createResponse() ([]byte, error) {
	response := newResponse()

	if r.winner != nil {
		response.Winner = r.winner.username
	}

	response.setScores(r.clients, r.nextTurn)

	return json.Marshal(response)
}

func (r *Room) handleError(client *Client, err error) {
	if client != nil {
		response, err := r.createErrResponse(err)
		if err != nil {
			log.Println("error marshalling scores", err)
			client.send <- []byte(err.Error())
		}
		client.send <- response
	}

}

func (r *Room) createErrResponse(err error) ([]byte, error) {
	response := &Response{}
	response.setErrMsg(err.Error())

	return json.Marshal(response)
}

func (r *Room) broadcast(message []byte) {
	for client := range r.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(r.clients, client)
		}
	}
}
