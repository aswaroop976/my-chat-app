package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust the origin checking for your requirements
	},
	// ReadBufferSize:  1024, //idk if I need these or not
	// WriteBufferSize: 1024,
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Infinite loop for reading messages from the WebSocket
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}
		fmt.Printf("message: %s", string(msg))
		if err = ws.WriteMessage(msgType, msg); err != nil {
			return
		}
		// var msg string
		// // Read in a new message as JSON and map it to a Message object
		// err := ws.ReadJSON(&msg)
		// if err != nil {
		// 	log.Printf("error: %v", err)
		// 	break
		// }
		// // Here, you can process the message or send it to other clients
		// log.Printf("Received: %s", msg)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("very crazy"))
	})

	http.HandleFunc("/ws", handleConnections)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
