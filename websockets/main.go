// Using go-gorilla websocket
// And rewriting their example chat client/hub
// Trying it in a single file for brevity sake
package main

import (
	"bytes" // for Client
	"flag"
	"fmt"
	"net/http"
	"text/template"
	"time" // for Client

	"github.com/gorilla/websocket" // for Client
)

// What does flag.String do?
var addr = flag.String("addr", ":8080", "http service address")

// Use index template
var index = template.Must(template.ParseFiles("index.html"))

/* Main Functions */
// This is the handler for the home page template
func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	// Only serve root index
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	// No posts!
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// cause bad things happen when we don't allow unicode
	index.Execute(w, r.Host)
}

func serveWebsockets(w http.ResponseWriter, r *http.Request) {
	// Wrapper around Hub's serve Websocket function
	serveWs(hub, w, r)
}

// Handle Functions and Server
func main() {
	flag.Parse() // Why? What?
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("ws", serveWebsockets)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Printf("Listen Error %s", err)
	}
}

/* Hub Struct and Functions
no imports? Uses Client Struct
*/
type Hub struct {
	// Registered Clients
	clients map[*Client]bool // true for is connected?
	// Messages from Clients
	// Messages are in bytes
	broadcast chan []byte
	// Register requests from Clients
	// Takes Client struct
	register chan *Client
	// Unregister, disconnect Clients
	unregister chan *Client
}

// Constructor for Hub
func newHub() *Hub {
	return &Hub{ // Return a pointer to new Hub
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		clients:    make(map[*Client]bool),
		unregister: make(chan *Client),
	}
}

// Single Hub method is for run()ning hub
func (h *Hub) run() {
	for {
		select {
		// incoming chan of register
		case client := <-h.register:
			// adds client to clients map
			h.clients[client] = true
		// Incoming chan of disconnects
		case client := <-h.unregister:
			// If client exists in map...
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- messae:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

/* Client Struct and method/functions
 */
type Client struct {
	hub *Hub // takes a hub duh
	// Websocket Connection, keep track of...
	conn *websocket.Conn
	// Buffered channel of outbound messages
	// In bytes, with 1 kb buffer?
	send chan []byte
	// What about incoming?
	// No constructor, rather constructed from Hub?
}
