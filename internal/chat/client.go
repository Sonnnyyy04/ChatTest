package chat

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func HandleMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}
		fmt.Print(msg)
	}
}

func SendMessage(conn net.Conn, message string) {
	_, err := fmt.Fprintf(conn, "%s\n", message)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func StartMessaging(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter message (or 'EXIT' to quit):\n")

	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "EXIT" {
			SendMessage(conn, message)
			break
		}

		SendMessage(conn, message)
	}
}
