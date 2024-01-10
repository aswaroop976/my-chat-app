package main

// to run type in "go run main.go hub.go"

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	// "github.com/rs/cors"

	// "text/template"

	// "github.com/go-sql-driver/mysql"

	"github.com/gorilla/websocket"
)

// idk I was doing some cloud-sql stuff here might continue later
//
//	func connect() {
//		cleanup, err := mysql.RegisterDriver("cloudsql-mysql")
//		if err != nil {
//			// ... handle error
//		}
//		// call cleanup when you're done with the database connection
//		defer cleanup()
//		var (
//			dbUser                 = "root"
//			dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
//			dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
//			instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
//			usePrivate             = os.Getenv("PRIVATE_IP")
//		)
//		dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
//			dbUser, dbPwd, dbName)
//		db, err := sql.Open(
//			"cloudsql-mysql",
//			"real-time-chat-app-409421:us-central1:fighting-game-main-db",
//		)
//		// ... etc
//	}

// struct used when users try to login
// used to store the username and password passed in as a json object
type Login struct {
	Username string `json:"username":`
	Password string `json:"password"`
}

// same idea with the login struct
// used to store the username, email, and password pased in when a user signs up
type Signup struct {
	Username string `json:"username":`
	Password string `json:"password"`
	Email    string `json: "email"`
}

// struct to write messages back to frontend in a json format
type Message struct {
	Content string `json:"content"`
}

// upgrades the http connection to a websocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust the origin checking for your requirements
	},
}

// run function for the hub(see definition in hub.go). Runs the hub struct, controls sending messages to the clients.
// Also handles registering and unregisterign of clients
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

// checks the send channel of clients, making sure to write messages to the websocket connecton once they appear on the send channel
func (c *Client) writePump() {
	for {
		message, ok := <-c.send
		if !ok {
			// The channel is closed, handle disconnection
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			// Handle errors (like disconnection)
			break
		}
	}
}

// like it sounds handles websocket connections when they initially connect, and listens for messages from these connections, and broadcasts any messages
// it receives.
func handleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	client := &Client{conn: ws, send: make(chan []byte, 256)}
	hub.register <- client
	go client.writePump()
	defer func() { hub.unregister <- client }()
	// Infinite loop for reading messages from the WebSocket
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		hub.broadcast <- message
	}
}

// handles sign in by parsing through POST request sent by the frontend with data regaring a new user trying to sign in
// if sign in is successful sends a message back to frontend indicating a page redirect
func handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	jsonFeed, err := io.ReadAll(r.Body)
	signup := Signup{}
	json.Unmarshal([]byte(jsonFeed), &signup)
	fmt.Println("Username", signup.Username, "Password", signup.Password, "Email", signup.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//need to add some checking to make sure that the username, email, and password are all valid
	msg := Message{Content: "Redirect"}
	json.NewEncoder(w).Encode(msg)

}

// handles log in by parsing through POST request sent by frontend with user trying to log back in.
// if login is successful sends a message back to frontend indicating a page redirect
func handleLogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	jsonFeed, err := io.ReadAll(r.Body)
	login := Login{}
	json.Unmarshal([]byte(jsonFeed), &login)
	fmt.Println("Username", login.Username, "Password", login.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if login.Username == "test" && login.Password == "test" { // replace with an actual check for username and passwd
		//render the frontend
		// http.ServeFile(w, r, filepath.Join("frontend", "index.html")) //I'm a fool for thinking this could work
		msg := Message{Content: "Allow login"}
		json.NewEncoder(w).Encode(msg)
	} else {
		msg := Message{Content: "No"}
		json.NewEncoder(w).Encode(msg)
	}
}

func main() {
	hub := newHub() // initializes a new hub struct
	go hub.run()    // runs the run() function in a seperate go routine
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("very crazy"))
	// })
	// fileServer := http.FileServer(http.Dir("./static"))
	// http.Handle("/", fileServer)
	fs := http.FileServer(http.Dir("frontend"))
	// frontend_fs := http.FileServer(http.Dir("frontend"))

	// Serve static files
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	// Custom handler for the root URL
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve index.html for the root URL
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("frontend", "login.html"))
			// return
		}
		// Serve static files for other URLs
		fs.ServeHTTP(w, r)
	})

	// http.Handle("/frontend/", http.StripPrefix("/frontend/", frontend_fs))
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		// Serve index.html for the root URL
		if r.URL.Path == "/chat" {
			http.ServeFile(w, r, filepath.Join("frontend", "index.html"))
			// return
		}
		// Serve static files for other URLs
		// frontend_fs.ServeHTTP(w, r)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogIn(w, r)
	})
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handleSignIn(w, r)
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
