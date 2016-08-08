package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	publicDirectory := "/lol"
	fmt.Println("Serving Files on port 1337")
	log.Fatal(http.ListenAndServe("0.0.0.0:1337", http.FileServer(http.Dir(publicDirectory))))

}
