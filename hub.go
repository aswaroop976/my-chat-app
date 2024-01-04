package main

import (
	"github.com/gorilla/websocket"
)

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

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte  // Broadcasts messages to the clients
	register   chan *Client // Registers new clients
	unregister chan *Client // Removes clients
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
