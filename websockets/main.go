// Using go-gorilla websocket
// And rewriting their example chat client/hub
package main

import (
	"flag"
	"fmt"
	"net/http"
	"text/template"
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

/* Hub Struct and Functions */