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
var usernames map[int]string

type Client struct {
	username string
	socketNumber int
	client net.Conn

}

func (client *Client) List() {
	fmt.Printf(client.username)
}

func main() {
	server, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if err != nil {
		panic("couldn't start listening: ")
	}
	conns := acceptClients(server)
	for {
		go acceptData(<-conns)

	}
}

//make connection
//accept new connections
func acceptClients(listener net.Listener) chan net.Conn {
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
			//username,_ :=  bufio.NewReader(client).ReadString('\n')
			//sockNum := client.RemoteAddr.String()[9:11]
			//username := "woosy"
			//sockNum := 10
			//c := Client{username: username, socketNumber: sockNum, client: client}
			//print info of user
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

//broadcast message

//forward message
func acceptData(client net.Conn) {
	fmt.Print("\n")
	for {
		msg, _ := bufio.NewReader(client).ReadString('\n')// what does the \n do?
		username := "<" + client.RemoteAddr().String()[9:11] + ">"
		if msg == ""{
			//fmt.Print(username, " leaves the chat")
			client.Close() // doesn't work
			continue
		}
		fmt.Print(username, string(msg))
	}
}
