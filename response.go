package main

// Response represents websocket message structure
type Response struct {
	Winner string
	Scores map[string]int
}

func newResponse() *Response {
	return &Response{
		Scores: make(map[string]int),
	}

}

func (r *Response) setScores(clients map[*Client]int) {
	if r.Scores == nil {
		r.Scores = make(map[string]int)
	}

	for client, score := range clients {
		r.Scores[client.username] = score
	}
}
