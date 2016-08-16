package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	publicDirectory := "/home/fen/access/Public"
	fmt.Println("Serving Files on port 80")
	log.Fatal(http.ListenAndServe("0.0.0.0:80", http.FileServer(http.Dir(publicDirectory))))

}
