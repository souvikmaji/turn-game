package main

import (
	"log"
)

// Lobby maintains the set of client rooms
type Lobby struct {
	// Register requests from the clients.
	register chan *Client

	// Active rooms with clients
	rooms map[*Room]bool
}

func newLobby() *Lobby {
	return &Lobby{
		register: make(chan *Client),
		rooms:    make(map[*Room]bool),
	}
}

func (l *Lobby) run() {
	for {
		select {
		case client := <-l.register:
			room := l.findRoom()

			// adding client to a room
			client.room = room
			client.room.register <- client

			// adding this new room in the lobby's room list
			l.rooms[room] = true
		}
	}
}

func (l *Lobby) findRoom() *Room {
	for room := range l.rooms {

		if !room.isFull {
			return room
		}
	}

	log.Println("creating new room")
	return newRoom()
}
