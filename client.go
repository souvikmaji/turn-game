package main

import (
	"github.com/bxcodec/faker/v3"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	outBufferSize = 256
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true }, // for external ws testers
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the lobby.
type Client struct {
	username string

	// helps to determine whose turn is next
	position int

	// room this client belongs to
	room *Room

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		username: faker.Username(), //TODO: remove faker dependency
		conn:     conn,
		send:     make(chan []byte, outBufferSize),
	}
}

// read pumps messages from the websocket connection to the lobby.
//
// The application runs read method in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) read() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v\n", err)
			}
			break
		}

		// ignore messages other than roll
		if string(message) == "roll" {
			c.room.playerMove <- c
		}

	}
}

// write pumps messages from the lobby to the websocket connection.
//
// A goroutine running write method is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) write() {

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Println("client closed the channel")
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		}
	}
}

// serveWebsocket handles new websocket requests from the peer.
func serveWsClient(lobby *Lobby, w http.ResponseWriter, r *http.Request) {
	// upgrade http connection to websocket protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := newClient(conn)

	lobby.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()
}
