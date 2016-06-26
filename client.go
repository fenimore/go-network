package main

import "net"
//import "bufio"
import "fmt"
import "bufio"
import "os"
import "strconv"
import "math/rand"
import "time"

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
	fmt.Println(conn, msg)
	conn.Write([]byte(msg + ">"))
//	status, _ := bufio.NewReader(conn).readString('\n')
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message : ")
		fmt.Print(strconv.Itoa(len(names)))
		//fmt.Println(strconv.Itoa(rand.Seed(23).Intn(len(names))))
		msg, _ := reader.ReadString('\n')
		conn.Write([]byte(msg + "\n"))
        }
}
