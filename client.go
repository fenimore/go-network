package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
)

// Client
func main() {
	conn, _ := net.Dial("tcp", "localhost:8998")
	welcomeUser()
	
	outgoing := make(chan string)
	incoming := make(chan string)
	var lastMessage string
	go checkIncoming(conn, incoming, lastMessage)
	go checkOutgoing(outgoing, &lastMessage)
	
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
func checkOutgoing(outs chan<- string, last *string) {
	inputReader := bufio.NewReader(os.Stdin)
	//var screen *bytes.Buffer = new(bytes.Buffer)
	//var output *bufio.Writer = bufio.NewWriter(os.Stdout)
	for {
		outgoing, err := inputReader.ReadString('\n')
		if err != nil {
			break
		}
		last = &outgoing
		//outs <- outgoing
		//output.WriteString("\033[2J")
		//fmt.Fprintf(screen, "\033[%d;%dH", 1, 1)
		outs <- outgoing
	}
}

// Check for incoming from server
func checkIncoming(conn net.Conn,
	ins chan<- string, last string) {

	connReader := bufio.NewReader(conn)

	for {
		incoming, _ := connReader.ReadString('\n')
		matched, _ := regexp.MatchString(last, incoming)
		if matched {
			fmt.Print("itss a match")
		}
		ins <- incoming
	}
}


func welcomeUser(){
	welcomeMessage := "Welcome"
	fmt.Println(welcomeMessage)
}

