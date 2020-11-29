package main

// Response represents websocket message structure
type Response struct {
	Winner  string
	Scores  map[string]int
	IsError bool
	ErrMsg  string
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

func (r *Response) setErrMsg(msg string) {
	r.IsError = true
	r.ErrMsg = msg
}
