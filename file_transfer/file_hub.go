package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// So whats going on with buffers?

const BUFFERSIZE = 1024 // About 1 KB at a time?
const PORT = 1337

func main() {

	// Listen on Port

	hub, err := net.Listen("tcp", "localhost:"+strconv.Itoa(PORT))
	if err != nil {
		fmt.Println("Error Listening: ", err)
		os.Exit(1)
	}

	defer hub.Close()
	fmt.Println("Listening on, " + strconv.Itoa(PORT))

	fmt.Print("Send file: ")
	// Choose file to send
	inputReader := bufio.NewReader(os.Stdin)
	filename, err := inputReader.ReadString('\n')
	filename = strings.TrimRight(filename, "\r\n")
	if err != nil {
		fmt.Println(err)
	}

	for {
		connection, err := hub.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Client Connected")
		go sendFile(connection, filename)
	}
}

// Fill message to buffer size
func fillString(result string, length int) string {
	for {
		resLength := len(result)
		if resLength < length {
			result = result + ":"
			continue
		}
		break
	}
	return result
}

func sendFile(conn net.Conn, name string) {
	defer conn.Close()
	// open the file that needs to be sent
	fmt.Println(name)
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Info has the file stats...
	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(info.Size(), 10), 10) // base 10
	fmt.Println(strconv.FormatInt(info.Size(), 10))
	fileName := fillString(info.Name(), 64) // A 64 possible size?
	fmt.Println("Sending filename and filesize")

	// Write first 10- bytes telling the filesize
	// Then write 64 bytes to vlient telling the the filename
	conn.Write([]byte(fileSize)) // Ten Bytes
	conn.Write([]byte(fileName)) // Sixty Four Bytes

	// The buffer, which will be a []byte
	// It is reused, and reused until the file has been totally read
	fmt.Println("Sending File")
	buf := make([]byte, BUFFERSIZE)

	// For loop until the file has totally been read/written
	// This is interesting. Using the same buffer, Read onto it
	// And then write using it onto the Connection, which is a Writer
	// I imagine.
	for {
		_, err = file.Read(buf) // The Reader takes a buffer size
		if err == io.EOF {      // it is the end of file
			break
		}
		conn.Write(buf)
	}
	fmt.Println("Sent!")
	return
}
