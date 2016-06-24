package main

import "net"
//import "bufio"
import "fmt"
import "bufio"
import "os"


func main() {
	conn, _ := net.Dial("tcp", "localhost:3540")
	var msg string
	msg = "yo"
	fmt.Fprintf(conn, msg)
	conn.Write([]byte(msg + "\n"))
//	status, _ := bufio.NewReader(conn).readString('\n')
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message : ")
		msg, _ := reader.ReadString('\n')
		conn.Write([]byte(msg + "\n"))
        }
}
