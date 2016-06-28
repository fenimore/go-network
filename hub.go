//http://synflood.at/tmp/golang-slides/mrmcd2012.html#9
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)


type Client struct {
	conn net.Conn
	ch chan<- string
}

func main(){
	port := "8998"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	messages := make(chan string)
	clients := make(chan Client)
	disconnects := make(chan net.Conn)

	// goroutine
	// Check for incoming messages, clients, or disconnects
	go inCome(messages, clients, disconnects)

	for {
		// check if channel is ready
		readySocket, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// goroutine for sock lifecycle
		go outGo(readySocket, messages, clients, disconnects)
	}
}

//Select from one of three channels:
//New Clients, New Messages, or New Disconnects
//Run in a goroutine
func inCome(msgs <-chan string, conns <-chan Client, discs <-chan net.Conn){
	connections := make(map[net.Conn]chan<- string)
	for {
		select {
		case msg := <- msgs:
			//print
			for _, val := range connections {
				go func(ch chan<- string){
					ch <- msg
				}(val)
			}
		case client := <-conns:
			fmt.Printf("New client: %v\n", client.conn)
			connections[client.conn] = client.ch
		case disc := <-discs:
			fmt.Printf("Client disconnects: %v\n", disc)
			delete(connections, disc)
		}
	}
}

//Outgoing messages.
func outGo(conn net.Conn, msgs chan<- string,
	conns chan<- Client, discs chan<- net.Conn) {
	channel := make(chan string)
	messages := make(chan string)
	conns <- Client{conn, channel}

	// as for username then
	// read input from client and add to messages
	go func() {
		defer close(messages)
		bufc := bufio.NewReader(conn)
		conn.Write([]byte("What is your name?\n"))
		nick, _, _ := bufc.ReadLine()
		handle := string(nick)
		conn.Write([]byte("Welcome!\n"))
		messages <- handle + " has joined chat"
		for {
			line, _, err := bufc.ReadLine()
			if err != nil {
				break
			} // line is still sent to self
			messages <- "<"+handle+"> "+string(line)
		}
		messages <- handle + " has left chat"
	}()// end of go func
Loop: // label loop for breaking out of it
	for {
		select {
		case msg, ok := <-messages:
			if !ok{
				break Loop
			}
			msgs <- msg
		case msg := <-channel: // don't understand
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				break Loop
			}
		}
	}

	conn.Close()
	fmt.Printf("Connection from %v, closed\n", conn.RemoteAddr())
	discs <- conn
}
