package main

import (
	"flag"
	"fmt"
	"log"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "the server address")
	flag.Parse()

	server := NewServer(*listenAddr)
	fmt.Println("server running on port:", *listenAddr)
	log.Fatal(server.Start())

	flag.Parse()
}
