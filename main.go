package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

// use default options
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // for external ws testers
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println("serving home")
	http.ServeFile(w, r, "home.html")
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("serving websocket")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()

	log.Println("starting server at: ", *addr)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWebsocket)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
