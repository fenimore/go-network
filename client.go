package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	// "strconv"
	"math/rand"
	"time"
)



func sendData(client net.Conn){
	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		msg, _ := reader.ReadString('\n')
		client.Write([]byte(msg + "\n"))
	}()
}

func receiveData(client net.Conn, data chan string) chan string{
	go func(){
		reader := bufio.NewReader(client)
	// flush?
		recv, _ := reader.ReadString('\n')
		fmt.Print(recv)
		data <- recv
	}()
	return data
}



func main() {
	//names := []string{"alpha","beta","gamma",}
	names:= make([]string, 0, 10)
	names= append(names,
		"alpha",
		"beta",
		"gamma",
		"zeta",
		"meta",
		"greta",
		"woops",
		"frinzipat",
		"calhou",
	)
	conn, _ := net.Dial("tcp", "localhost:3540")
	var msg string
	rand.Seed(int64(time.Now().Nanosecond()))
	msg = names[rand.Intn(len(names))]
	data := make(chan string)
	//fmt.Println(conn, msg)
	conn.Write([]byte(msg + ">"))
	data = receiveData(conn, data)
	fmt.Print(<-data)
	//
	for {
		sendData(conn)
		incoming := receiveData(conn, data)
		fmt.Print(<-incoming)
	}
}
