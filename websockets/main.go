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

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second
	// Time allowed to read the next pong messsage from peer
	pongWait = 60 * time.Second
	//Send pings to peer, must be less than ping wait
	pingPeriod = (pongWait * 9) / 10
	// Maximun message size allowed from peer.
	maxMessageSize = 512 // 512 what?
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// WebsocketUpgrader... does something?
var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024 // one kb
	WriteBufferSize: 1024 // must they be the same?
}

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
// Run() is an infinite loop taking in the different
// Hub channels, and performing either:
// connects, disconnects, or broadcasts
// depending on these channels
// This method must be run in a goroutine
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
// Client is a middleman between
// websocket connection and the hub
// It keeps track of the connection, the hub
// and the outgoing messages in a chan
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
// write method writes a message 
func (c *Client) write(mt int, payloud []byte) error {
	// I don't understand this method...
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, payload)
}
// readPump pumps message from the websocket connection to the hub
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// Hub closed channel.
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// TOo tired to continue
		}
	}
			
}
// writePump pumps messages from the hub to the connection

// serveWs handles websocket requests from peer
// The handler for the /ws route

//
