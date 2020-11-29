package main

import (
	"log"
)

// Response represents websocket message structure
type Response struct {
	Winner   string
	NextMove string
	Scores   map[string]int
	IsError  bool
	ErrMsg   string
}

func newResponse() *Response {
	return &Response{
		Scores: make(map[string]int),
	}

}

func (r *Response) setScores(clients map[*Client]int, nextMove int) {
	if r.Scores == nil {
		r.Scores = make(map[string]int)
	}

	for client, score := range clients {
		if client.position == nextMove {
			r.NextMove = client.username
		}
		r.Scores[client.username] = score

	}
}

func (r *Response) setErrMsg(msg string) {
	r.IsError = true
	r.ErrMsg = msg
}
