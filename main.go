package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var room *Room

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println("serving home")
	http.ServeFile(w, r, "home.html")
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("serving websocket")
	serveWsClient(room, w, r)
}

func main() {
	flag.Parse()

	log.Println("starting server at: ", *addr)

	room = newRoom()
	go room.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWebsocket)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
