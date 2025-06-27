package main

import (
	"chat_test/internal/chat"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	username, password := chat.GetCredentials()
	chat.SendCredentials(conn, username, password)
	if !chat.Login(conn) {
		return
	}

	go chat.HandleMessages(conn)

	chat.StartMessaging(conn)
}
