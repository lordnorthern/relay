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

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/relay.html")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<html>`)
}

func relay(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]interface{})
	r.ParseMultipartForm(10 << 20)
	err := json.NewDecoder(r.Body).Decode(&body)
	var b []byte
	if err != nil {
		contact := make(map[string]string)
		for k, v := range r.Form {
			contact[k] = v[0]
		}
		b, err = json.Marshal(contact)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		b, err = json.Marshal(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	headers, _ := json.Marshal(r.Header)

	cont := Contact{Header: string(headers), Body: string(b)}
	b, err = json.Marshal(cont)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for client := range connectedClients {
		if err := client.WriteMessage(1, b); err != nil {
			log.Println(err)
			client.Close()
			delete(connectedClients, client)
			fmt.Println("Client Disconnected")
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

	http.HandleFunc("/home", home)
	http.HandleFunc("/relay", relay)
	http.HandleFunc("/ws", wsEndpoint)
	fmt.Println("Listening on port " + port)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":"+port, nil)
	fmt.Println(err)
}

type Contact struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}
