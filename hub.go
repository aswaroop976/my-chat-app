package main

import (
	"github.com/gorilla/websocket"
)

// made when a user makes a websocket connection to the server.
// contains a websocket connection to access the connection
// and contains the send channel, this channel contains messages to be written to the websocket connection
type Client struct {
	conn *websocket.Conn
	// userId   int
	// username string
	// email    string
	// password string
	send chan []byte
}

// use channels here b/c we want to work with go-routines(similar to threads), channels allow data to be shared
// across go-routines without the need for locks. Synchronization is automatic, this works because of the blocking behavior of channels

// struct used to control all the client websocket connections to the server, and to broadcast messages
// clients attribute is a map indicating whether or not a client is connected(true) or disconnected(false)
// broadcast channel used to broadcast messages to all the clients
// register and unregister channels used to register and remove clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte  // Broadcasts messages to the clients
	register   chan *Client // Registers new clients
	unregister chan *Client // Removes clients
}

// initializer function for the hub struct
func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
