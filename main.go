package main

import (
	"log"
	"net/http"

	// "github.com/go-sql-driver/mysql"

	"github.com/gorilla/websocket"
)

// idk I was doing some cloud-sql stuff here might continue later
// func connect() {
// 	cleanup, err := mysql.RegisterDriver("cloudsql-mysql")
// 	if err != nil {
// 		// ... handle error
// 	}
// 	// call cleanup when you're done with the database connection
// 	defer cleanup()
// 	var (
// 		dbUser                 = "root"
// 		dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
// 		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
// 		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
// 		usePrivate             = os.Getenv("PRIVATE_IP")
// 	)
// 	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
// 		dbUser, dbPwd, dbName)
// 	db, err := sql.Open(
// 		"cloudsql-mysql",
// 		"real-time-chat-app-409421:us-central1:fighting-game-main-db",
// 	)
// 	// ... etc
// }

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust the origin checking for your requirements
	},
	// ReadBufferSize:  1024, //idk if I need these or not
	// WriteBufferSize: 1024,
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func handleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	client := &Client{conn: ws, send: make(chan []byte, 256)}
	hub.register <- client

	defer func() { hub.unregister <- client }()
	// Infinite loop for reading messages from the WebSocket
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		hub.broadcast <- message
		// msgType, msg, err := ws.ReadMessage()
		// if err != nil {
		// 	return
		// }
		// fmt.Printf("message: %s", string(msg))
		// if err = ws.WriteMessage(msgType, msg); err != nil {
		// 	return
		// }
		// // var msg string
		// // // Read in a new message as JSON and map it to a Message object
		// // err := ws.ReadJSON(&msg)
		// // if err != nil {
		// // 	log.Printf("error: %v", err)
		// // 	break
		// // }
		// // // Here, you can process the message or send it to other clients
		// // log.Printf("Received: %s", msg)
	}
}

func main() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("very crazy"))
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, w, r)
	})

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
