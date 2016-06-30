package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// Client
func main() {
	conn, _ := net.Dial("tcp", "localhost:8998")
	welcomeUser()
	
	outgoing := make(chan string)
	incoming := make(chan string)

	go checkIncoming(conn, incoming)
	go checkOutgoing(outgoing)
	
	for {
		select {
		case msg := <-outgoing:
			fmt.Fprintf(conn, msg)
			break
		case msg := <-incoming:
			fmt.Print("\r")
			fmt.Print(msg) // the server sends \n
			fmt.Print("> ")
			break
		}
	}
}

// Check for User input for outgoing channel
func checkOutgoing(outs chan<- string) {
	inputReader := bufio.NewReader(os.Stdin)
	//var screen *bytes.Buffer = new(bytes.Buffer)
	//var output *bufio.Writer = bufio.NewWriter(os.Stdout)
	for {
		outgoing, err := inputReader.ReadString('\n')
		if err != nil {
			break
		}
		outs <- outgoing
	}
}

// Check for incoming from server
func checkIncoming(conn net.Conn,
	ins chan<- string) {

	connReader := bufio.NewReader(conn)

	for {
		incoming, _ := connReader.ReadString('\n')
		ins <- incoming
	}
}


func welcomeUser(){
	welcomeMessage := `
    ___       ___       ___       ___   
   /\  \     /\  \     /\__\     /\  \  
  /::\  \   /::\  \   /:/__/_   /::\  \ 
 /::\:\__\ /:/\:\__\ /::\/\__\ /:/\:\__\
 \:\:\/  / \:\ \/__/ \/\::/  / \:\/:/  /
  \:\/  /   \:\__\     /:/  /   \::/  / 
   \/__/     \/__/     \/__/     \/__/  
`
	fmt.Println(welcomeMessage)
}

