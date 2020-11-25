package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var lobby *Lobby

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println("serving home")
	http.ServeFile(w, r, "home.html")
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("serving websocket")
	serveWsClient(lobby, w, r)
}

func main() {
	flag.Parse()

	log.Println("starting server at: ", *addr)

	lobby = newLobby()
	go lobby.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWebsocket)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
