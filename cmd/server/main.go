package main

import (
	"chat_test/internal/chat"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Println("Server started on localhost:8081")

	chat.StartServer(listener)
}
