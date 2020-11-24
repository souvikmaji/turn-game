package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println("serving home")
	http.ServeFile(w, r, "home.html")
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("serving websocket")
	fmt.Fprintf(w, "websocket")
}

func main() {
	flag.Parse()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWebsocket)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
