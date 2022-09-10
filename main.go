package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const port = "7717"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connectedClients = make(map[*websocket.Conn]bool)

func relay(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for client := range connectedClients {
		if err := client.WriteMessage(1, b); err != nil {
			log.Println(err)
			client.Close()
			delete(connectedClients, client)
		}
	}
	fmt.Fprintf(w, "OK\n")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	connectedClients[ws] = true
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			delete(connectedClients, conn)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			delete(connectedClients, conn)
			return
		}
	}
}

func main() {

	http.HandleFunc("/relay", relay)
	http.HandleFunc("/ws", wsEndpoint)
	fmt.Println("Listening on port " + port)
	http.ListenAndServe(":"+port, nil)
}
