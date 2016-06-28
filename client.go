package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:8998")
	outgoing := make(chan string)
	incoming := make(chan string)
	go checkIncoming(conn, incoming)
	go checkOutgoing(outgoing)
	for {
		select {
		case msg := <-outgoing:
			fmt.Fprintf(conn, msg)// don't \n
			break
		case msg := <-incoming:
			fmt.Print(msg)// the server sends \n
			break
		}
	}
}

func checkOutgoing(outs chan<- string){
	inputReader := bufio.NewReader(os.Stdin)
	for{
		outgoing, _ := inputReader.ReadString('\n')
		outs <- outgoing
	}
}

func checkIncoming(conn net.Conn, ins chan<- string){
	connReader := bufio.NewReader(conn)
	for{
		incoming, _ := connReader.ReadString('\n')
		ins <- incoming
	}
}
