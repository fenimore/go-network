// http://www.mrwaggel.be/post/golang-sending-a-file-over-tcp/
// ^^ Adapted from that

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

const BUFFERSIZE = 1024
const PORT = 1337

func main() {
	// Doesn't seem to connect when I enter my internall addres? :S
	fmt.Print("Enter internal IP (or leave blank for localhost): ")
	inputReader := bufio.NewReader(os.Stdin)
	IP, err := inputReader.ReadString('\n')
	IP = strings.TrimRight(IP, "\r\n")
	if IP == "" {
		IP = "localhost"
	}

	conn, err := net.Dial("tcp", IP+":"+strconv.Itoa(PORT))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected, start receiving, first size then name then actual data")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	conn.Read(bufferFileSize) // first read size
	fmt.Printf("File size: %s\n", strings.Trim(string(bufferFileSize), ":"))
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	conn.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	fmt.Printf("File Name: %s\n", fileName)
	// Create the file to write
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	// Why int64?
	var totalRecv int64 // Total data we have received
	for {
		// finished receiving file because
		// if what is left to read, is less than the buffersize
		// then the file has been read?
		if (fileSize - totalRecv) < BUFFERSIZE {
			// I don't understand...
			io.CopyN(file, conn, (fileSize - totalRecv))
			// empty the remaining bytes that we don't need from the network buffer
			// Read onto a new buffer?
			conn.Read(make([]byte, (totalRecv+BUFFERSIZE)-fileSize))
			// Finished reading file
			break
		}
		io.CopyN(file, conn, BUFFERSIZE)
		// Increment Counter
		totalRecv += BUFFERSIZE
	}
	fmt.Println("Received file")
}
