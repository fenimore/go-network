// chat server
package main

import (
	"bufio"
	"fmt"
	"net"
)

type ChatRoom struct {
	listener net.Listener
	name string
	port string
	usernames map[string]string
	//socket net.Conn
	clients []net.Conn
	outgoing chan string
}

func (room *ChatRoom) launchServer(){
	// listen
	fmt.Println(room.name)
	server, err := net.Listen("tcp", ":" + room.port)
	if err != nil {
		panic("couldn't start listening: ")
	}
	room.listener = server
	//room.outgoing := make(chan string)
	//conns := room.acceptConnections()
	for {
		conns := room.acceptConnections()
		outgoings := make(chan string)
		go room.acceptData(<-conns, outgoings)
		go room.echoData(<-outgoings)
		//close(outgoings)
	}
}

func (room *ChatRoom) acceptConnections() chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func () {
		for {
			client, _ := room.listener.Accept()
			room.clients = append(room.clients, client)
			username, _ := bufio.NewReader(client).ReadString('>')
			fmt.Printf(username, " enters room")
			connectionKey := client.RemoteAddr().String()[9:11]
			room.usernames[connectionKey] = username
			client.Write([]byte(username + " logs in\n"))
			i++
			ch <- client
		}
	}()
	return ch
}

func (room *ChatRoom) acceptData(client net.Conn, in chan string){ // return outs?
	fmt.Print("\n")
	//out := make(chan string)
	for {
		incoming, _ := bufio.NewReader(client).ReadString('\n')
		id := client.RemoteAddr().String()[9:11]
		u, _ := room.usernames[id]
		username := "<" + u
		if incoming == "" {
			client.Close()
			//continue
		}
		message := username + " " + string(incoming)
		fmt.Print(username," ", string(incoming))
		in <- message
		//room.outgoing <- message
		//fmt.Print(<-room.outgoing)
		//room.echoData(message)
	}
}


func (room *ChatRoom) echoData(out string){
	message := out
	for _, client := range room.clients {
		client.Write([]byte("something has been sent\n"))
		client.Write([]byte(message))
		fmt.Println(client.RemoteAddr().String())
		//client.Close()
	}
	
		
}

	
func main(){
	//new chatroom
	//declare properties
	// launch server
	chatroom := new(ChatRoom)
	chatroom.name = "lobby"
	chatroom.port = "3540"
	chatroom.clients = make([]net.Conn, 0, 10)
	chatroom.usernames = make(map[string]string)
	chatroom.launchServer()
}


