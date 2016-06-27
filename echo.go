// $ 6g echo.go && 6l -o echo echo.6
// $ ./echo
//
//  ~ in another terminal ~
//
// $ nc localhost 3540

package main

import (
	"bufio"
	"fmt"
//	"io"
	"net"
	"strconv"
)

const PORT = 3540

type Client struct {
	handle string
	socketNumber int
	socket net.Conn

}

func (client *Client) List() {
	fmt.Printf(client.handle)
}

func main() {
	//var usernames map[int]string
	//Clients := make([]net.Conn, 10)
	var clients []*Client
	usernames := make(map[int]string)
	server, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if err != nil {
		panic("couldn't start listening: ")
	}
	conns := acceptClients(server, usernames, clients)
	for {
		go acceptData(<-conns, usernames)
		//fmt.Printf("%v", clients)

	}
}

//make connection
//accept new connections
func acceptClients(listener net.Listener, usernames map[int]string, clients []*Client) chan net.Conn {
	ch := make(chan net.Conn)
	//connections := make(chan Client)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if err != nil {
				fmt.Printf("couldn't accept connection")
				continue
			}
			username,_ :=  bufio.NewReader(client).ReadString('>')
			c := Client{handle: username, socketNumber: 10, socket: client,}
			clients = append(clients, &c)
			fmt.Printf(username)
			connKey, _ := strconv.Atoi(client.RemoteAddr().String()[9:11])
			usernames[connKey] = username
			// print all usernames????
			
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			fmt.Printf("%v",clients)
			ch <- client
		}
	}()
	return ch
}

//broadcast message
// Find access to all clients!??!?
func broadCast(message string, client net.Conn){
	client.Write([]byte(message + "\n"))
}

//forward message
func acceptData(client net.Conn, usernames map[int]string) {
	fmt.Print("\n")
	for {
		msg, _ := bufio.NewReader(client).ReadString('\n')// what does the \n do?
		id, _ := strconv.Atoi(client.RemoteAddr().String()[9:11])
		u, _ := usernames[id]
		username := "<" + u
		if msg == ""{
			//fmt.Print(username, " leaves the chat")
			client.Close() // doesn't work
			continue
		}
		fmt.Print(username, " ", string(msg))
		broadCast(string(msg), client)
	}
}
